package service

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"photoset/internal/repository"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// maxXDBFileSize 下载文件大小上限（200MB）
const maxXDBFileSize = 200 * 1024 * 1024

// allowedDownloadURLs 合法的 ip2region 下载地址白名单（精确匹配，防止指向攻击者自己的仓库）
var allowedDownloadURLs = map[string]bool{
	// ip2region v4
	"https://gitee.com/lionsoul/ip2region/raw/master/data/ip2region_v4.xdb":  true,
	"https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v4.xdb": true,
	// ip2region v6
	"https://gitee.com/lionsoul/ip2region/raw/master/data/ip2region_v6.xdb":  true,
	"https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v6.xdb": true,
}

// validateDownloadURL 校验下载 URL（多层防护：协议 + URL 路径白名单）
func validateDownloadURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("无效的URL格式: %v", err)
	}

	// 第 1 层：只允许 https（拒绝 http 明文、file、gopher、dict 等危险协议）
	if u.Scheme != "https" {
		return fmt.Errorf("不允许的协议: %s，只允许 https", u.Scheme)
	}

	// 第 2 层：URL 路径必须在白名单中（精确匹配，防止指向攻击者自己的仓库或恶意文件）
	normalized := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
	if !allowedDownloadURLs[normalized] {
		return fmt.Errorf("不允许的下载地址: %s", normalized)
	}

	return nil
}

// validateDownloadedFile 校验下载后的文件是合法的 xdb 数据库（防止下载到恶意文件）
func validateDownloadedFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer f.Close()

	// 第 3 层：文件大小限制
	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("无法获取文件信息: %v", err)
	}
	if stat.Size() > maxXDBFileSize {
		return fmt.Errorf("文件过大: %d bytes (上限 %d)", stat.Size(), maxXDBFileSize)
	}
	if stat.Size() < 1024 { // xdb 至少有几 KB
		return fmt.Errorf("文件过小，不是合法的数据库文件: %d bytes", stat.Size())
	}

	// 第 4 层：读取文件头，尝试验证是合法 xdb 格式
	// ip2region xdb 文件头包含 headerLength (2 bytes little-endian)
	header := make([]byte, 2)
	if _, err := io.ReadFull(f, header); err != nil {
		return fmt.Errorf("读取文件头失败: %v", err)
	}
	headerLen := binary.LittleEndian.Uint16(header)
	// header 长度通常在 256~65535 之间，不合理的值说明不是合法 xdb 文件
	if headerLen < 256 || headerLen > 65535 {
		return fmt.Errorf("文件格式不正确: header_length=%d", headerLen)
	}

	return nil
}

// IPGeoService IP地理位置服务
type IPGeoService struct {
	searcherV4  *xdb.Searcher
	searcherV6  *xdb.Searcher
	dbPathV4    string
	dbPathV6    string
	mu          sync.RWMutex
	configRepo  *repository.SiteSettingRepository
	stopChan    chan struct{}
}

// IPGeoConfig IP地理位置配置
type IPGeoConfig struct {
	Enabled         bool   `json:"enabled"`
	Mirror          string `json:"mirror"` // "gitee" 或 "github"
	DownloadURLV4   string `json:"download_url_v4"`
	DownloadURLV6   string `json:"download_url_v6"`
	UpdateDays      int    `json:"update_days"`
	LastUpdate      string `json:"last_update"`
	DatabasePathV4  string `json:"database_path_v4"`
	DatabasePathV6  string `json:"database_path_v6"`
}

// 镜像源地址映射
var mirrorURLs = map[string][2]string{
	"gitee": {
		"https://gitee.com/lionsoul/ip2region/raw/master/data/ip2region_v4.xdb",
		"https://gitee.com/lionsoul/ip2region/raw/master/data/ip2region_v6.xdb",
	},
	"github": {
		"https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v4.xdb",
		"https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v6.xdb",
	},
}

// NewIPGeoService 创建IP地理位置服务
func NewIPGeoService() *IPGeoService {
	service := &IPGeoService{
		configRepo: repository.NewSiteSettingRepository(),
		stopChan:   make(chan struct{}),
	}

	// 加载配置
	config := service.loadConfig()
	service.dbPathV4 = config.DatabasePathV4
	service.dbPathV6 = config.DatabasePathV6

	// 加载数据库
	if config.Enabled {
		if err := service.loadDatabase(); err != nil {
			log.Printf("[IPGeo] 加载数据库失败: %v, 尝试自动下载...", err)
			// 数据库文件不存在时，自动下载
			if err := service.UpdateDatabase(); err != nil {
				log.Printf("[IPGeo] 自动下载数据库失败: %v", err)
			}
		} else {
			log.Printf("[IPGeo] 数据库加载成功")
		}
	}

	// 启动自动更新检查
	go service.startAutoUpdate()

	return service
}

// loadConfig 加载配置
func (s *IPGeoService) loadConfig() IPGeoConfig {
	config := IPGeoConfig{
		Enabled:        true,
		Mirror:         "gitee",
		DownloadURLV4:  mirrorURLs["gitee"][0],
		DownloadURLV6:  mirrorURLs["gitee"][1],
		UpdateDays:     10,
		DatabasePathV4: "data/ip2region_v4.xdb",
		DatabasePathV6: "data/ip2region_v6.xdb",
	}

	// 从数据库加载配置
	settings, err := s.configRepo.GetAll()
	if err != nil {
		log.Printf("[IPGeo] 加载配置失败: %v", err)
		return config
	}

	if v, ok := settings["ip_geo_enabled"]; ok {
		config.Enabled = v == "true"
	}
	if v, ok := settings["ip_geo_mirror"]; ok && v != "" {
		config.Mirror = v
		// 根据 mirror 更新默认下载地址
		if urls, ok := mirrorURLs[v]; ok {
			config.DownloadURLV4 = urls[0]
			config.DownloadURLV6 = urls[1]
		}
	}
	if v, ok := settings["ip_geo_download_url_v4"]; ok && v != "" {
		config.DownloadURLV4 = v
	}
	if v, ok := settings["ip_geo_download_url_v6"]; ok && v != "" {
		config.DownloadURLV6 = v
	}
	if v, ok := settings["ip_geo_update_days"]; ok {
		fmt.Sscanf(v, "%d", &config.UpdateDays)
	}
	if v, ok := settings["ip_geo_last_update"]; ok {
		config.LastUpdate = v
	}
	if v, ok := settings["ip_geo_database_path_v4"]; ok && v != "" {
		config.DatabasePathV4 = v
	}
	if v, ok := settings["ip_geo_database_path_v6"]; ok && v != "" {
		config.DatabasePathV6 = v
	}

	return config
}

// saveConfig 保存配置
func (s *IPGeoService) saveConfig(config IPGeoConfig) error {
	settings := map[string]interface{}{
		"ip_geo_enabled":          fmt.Sprintf("%v", config.Enabled),
		"ip_geo_mirror":           config.Mirror,
		"ip_geo_download_url_v4":  config.DownloadURLV4,
		"ip_geo_download_url_v6":  config.DownloadURLV6,
		"ip_geo_update_days":      fmt.Sprintf("%d", config.UpdateDays),
		"ip_geo_last_update":      config.LastUpdate,
		"ip_geo_database_path_v4": config.DatabasePathV4,
		"ip_geo_database_path_v6": config.DatabasePathV6,
	}
	return s.configRepo.BatchUpsert(settings)
}

// loadDatabase 加载数据库到内存
func (s *IPGeoService) loadDatabase() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 加载 IPv4 数据库
	if err := s.loadSearcherV4(); err != nil {
		return fmt.Errorf("加载IPv4数据库失败: %v", err)
	}

	// 加载 IPv6 数据库（可选）
	if err := s.loadSearcherV6(); err != nil {
		log.Printf("[IPGeo] IPv6数据库加载失败（可选）: %v", err)
	}

	return nil
}

// loadSearcherV4 加载 IPv4 searcher
func (s *IPGeoService) loadSearcherV4() error {
	if _, err := os.Stat(s.dbPathV4); os.IsNotExist(err) {
		return fmt.Errorf("IPv4数据库文件不存在: %s", s.dbPathV4)
	}

	cBuff, err := os.ReadFile(s.dbPathV4)
	if err != nil {
		return fmt.Errorf("读取IPv4数据库文件失败: %v", err)
	}

	searcher, err := xdb.NewWithBuffer(xdb.IPv4, cBuff)
	if err != nil {
		return fmt.Errorf("创建IPv4查询对象失败: %v", err)
	}

	if s.searcherV4 != nil {
		s.searcherV4.Close()
	}

	s.searcherV4 = searcher
	return nil
}

// loadSearcherV6 加载 IPv6 searcher
func (s *IPGeoService) loadSearcherV6() error {
	if _, err := os.Stat(s.dbPathV6); os.IsNotExist(err) {
		return fmt.Errorf("IPv6数据库文件不存在: %s", s.dbPathV6)
	}

	cBuff, err := os.ReadFile(s.dbPathV6)
	if err != nil {
		return fmt.Errorf("读取IPv6数据库文件失败: %v", err)
	}

	searcher, err := xdb.NewWithBuffer(xdb.IPv6, cBuff)
	if err != nil {
		return fmt.Errorf("创建IPv6查询对象失败: %v", err)
	}

	if s.searcherV6 != nil {
		s.searcherV6.Close()
	}

	s.searcherV6 = searcher
	return nil
}

// getSearcher 根据 IP 版本获取对应的 searcher
func (s *IPGeoService) getSearcher(ip string) *xdb.Searcher {
	ver, err := xdb.VersionFromIP(ip)
	if err != nil {
		log.Printf("[IPGeo] 判断IP版本失败: %s, error: %v", ip, err)
		return nil
	}

	if ver == xdb.IPv6 {
		return s.searcherV6
	}
	return s.searcherV4
}

// 中国省级行政区白名单（34个）
var chinaProvinces = map[string]bool{
	"北京": true, "上海": true, "天津": true, "重庆": true,
	"河北": true, "山西": true, "辽宁": true, "吉林": true,
	"黑龙江": true, "江苏": true, "浙江": true, "安徽": true,
	"福建": true, "江西": true, "山东": true, "河南": true,
	"湖北": true, "湖南": true, "广东": true, "海南": true,
	"四川": true, "贵州": true, "云南": true, "陕西": true,
	"甘肃": true, "青海": true, "台湾": true,
	"内蒙古": true, "广西": true, "西藏": true, "宁夏": true, "新疆": true,
	"香港": true, "澳门": true,
}

// knownNonProvinceCities 已知会被 ip2region 误当作省份返回的地级市 -> 省份映射
var knownNonProvinceCities = map[string]string{
	"辽源": "吉林", "延边": "吉林", "通化": "吉林", "白山": "吉林",
	"白城": "吉林", "松原": "吉林", "四平": "吉林",
	"鞍山": "辽宁", "抚顺": "辽宁", "本溪": "辽宁", "丹东": "辽宁",
	"锦州": "辽宁", "营口": "辽宁", "阜新": "辽宁",
	"辽阳": "辽宁", "盘锦": "辽宁", "铁岭": "辽宁",
	"朝阳": "辽宁", "葫芦岛": "辽宁",
	"玉林": "广西", "桂林": "广西", "梧州": "广西",
	"北海": "广西", "防城港": "广西", "钦州": "广西",
	"贵港": "广西", "百色": "广西", "贺州": "广西",
	"河池": "广西", "来宾": "广西", "崇左": "广西",
	"中山": "广东", "珠海": "广东", "东莞": "广东", "佛山": "广东",
	"惠州": "广东", "江门": "广东", "揭阳": "广东",
	"茂名": "广东", "梅州": "广东", "清远": "广东",
	"汕头": "广东", "汕尾": "广东", "韶关": "广东",
	"深圳": "广东", "阳江": "广东", "云浮": "广东",
	"湛江": "广东", "肇庆": "广东", "潮州": "广东",
	"河源": "广东",
	"宁德": "福建", "莆田": "福建", "泉州": "福建",
	"漳州": "福建", "南平": "福建", "龙岩": "福建",
	"三明": "福建",
	"荆州": "湖北", "宜昌": "湖北", "襄阳": "湖北",
	"黄石": "湖北", "十堰": "湖北", "孝感": "湖北",
	"黄冈": "湖北", "咸宁": "湖北", "恩施": "湖北", "随州": "湖北",
	"吉安": "江西", "九江": "江西", "萍乡": "江西",
	"鹰潭": "江西", "赣州": "江西", "宜春": "江西",
	"上饶": "江西", "景德镇": "江西", "新余": "江西", "抚州": "江西",
	"芜湖": "安徽", "蚌埠": "安徽", "淮南": "安徽",
	"马鞍山": "安徽", "淮北": "安徽", "铜陵": "安徽",
	"安庆": "安徽", "黄山": "安徽", "滁州": "安徽",
	"阜阳": "安徽", "宿州": "安徽", "六安": "安徽",
	"亳州": "安徽", "池州": "安徽", "宣城": "安徽",
	"常州": "江苏", "苏州": "江苏", "无锡": "江苏",
	"徐州": "江苏", "南通": "江苏", "扬州": "江苏",
	"盐城": "江苏", "镇江": "江苏", "泰州": "江苏",
	"淮安": "江苏", "连云港": "江苏",
	"济南": "山东", "青岛": "山东", "淄博": "山东",
	"枣庄": "山东", "东营": "山东", "烟台": "山东",
	"潍坊": "山东", "济宁": "山东", "泰安": "山东",
	"威海": "山东", "日照": "山东", "临沂": "山东",
	"德州": "山东", "聊城": "山东", "滨州": "山东", "菏泽": "山东",
	"保定": "河北", "邯郸": "河北", "唐山": "河北",
	"秦皇岛": "河北", "邢台": "河北", "沧州": "河北",
	"承德": "河北", "张家口": "河北", "衡水": "河北", "廊坊": "河北",
	"洛阳": "河南", "开封": "河南", "安阳": "河南",
	"鹤壁": "河南", "新乡": "河南", "焦作": "河南",
	"濮阳": "河南", "许昌": "河南", "漯河": "河南",
	"三门峡": "河南", "南阳": "河南", "商丘": "河南",
	"信阳": "河南", "周口": "河南", "驻马店": "河南", "平顶山": "河南",
	"绍兴": "浙江", "嘉兴": "浙江", "湖州": "浙江",
	"金华": "浙江", "衢州": "浙江", "舟山": "浙江",
	"台州": "浙江", "丽水": "浙江", "温州": "浙江", "宁波": "浙江",
	"绵阳": "四川", "德阳": "四川", "宜宾": "四川",
	"南充": "四川", "广元": "四川", "资阳": "四川",
	"眉山": "四川", "泸州": "四川", "达州": "四川",
	"乐山": "四川", "凉山": "四川", "甘孜": "四川", "阿坝": "四川",
	"株洲": "湖南", "湘潭": "湖南", "衡阳": "湖南",
	"邵阳": "湖南", "岳阳": "湖南", "常德": "湖南",
	"张家界": "湖南", "益阳": "湖南", "郴州": "湖南",
	"永州": "湖南", "怀化": "湖南", "娄底": "湖南", "湘西": "湖南",
	"遵义": "贵州", "安顺": "贵州", "铜仁": "贵州",
	"黔东南": "贵州", "黔南": "贵州", "黔西南": "贵州",
	"毕节": "贵州", "六盘水": "贵州",
	"大同": "山西", "阳泉": "山西", "长治": "山西",
	"晋城": "山西", "朔州": "山西", "晋中": "山西",
	"运城": "山西", "忻州": "山西", "临汾": "山西", "吕梁": "山西",
	"宝鸡": "陕西", "咸阳": "陕西", "铜川": "陕西",
	"渭南": "陕西", "延安": "陕西", "汉中": "陕西",
	"榆林": "陕西", "安康": "陕西", "商洛": "陕西",
	"大理": "云南", "红河": "云南", "曲靖": "云南",
	"玉溪": "云南", "楚雄": "云南", "丽江": "云南",
	"文山": "云南", "普洱": "云南", "昭通": "云南",
	"保山": "云南", "德宏": "云南", "西双版纳": "云南",
	"张掖": "甘肃", "武威": "甘肃", "酒泉": "甘肃",
	"平凉": "甘肃", "庆阳": "甘肃", "天水": "甘肃",
	"定西": "甘肃", "嘉峪关": "甘肃", "金昌": "甘肃",
	"陇南": "甘肃", "临夏": "甘肃", "甘南": "甘肃",
	"海东": "青海", "海西": "青海", "海北": "青海",
	"黄南": "青海", "果洛": "青海", "玉树": "青海",
	"固原": "宁夏", "吴忠": "宁夏", "中卫": "宁夏", "石嘴山": "宁夏",
	"伊犁": "新疆", "哈密": "新疆", "喀什": "新疆",
	"阿克苏": "新疆", "和田": "新疆", "吐鲁番": "新疆",
	"塔城": "新疆", "阿勒泰": "新疆", "昌吉": "新疆",
	"克孜勒苏": "新疆", "巴音郭楞": "新疆", "博尔塔拉": "新疆",
}

// normalizeRegionName 去掉行政区划后缀（省/市/自治区/特别行政区）
// 用于统一映射表 key 的匹配
func normalizeRegionName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" || name == "0" {
		return name
	}
	// 按优先级去掉后缀（长后缀先匹配）
	for _, suffix := range []string{"特别行政区", "自治区", "省", "市"} {
		if strings.HasSuffix(name, suffix) {
			return strings.TrimSuffix(name, suffix)
		}
	}
	return name
}

// GetLocation 获取IP地理位置（省份级别）
// 对国外IP返回国家名，国内IP返回省份名（确保不出现城市级）
func (s *IPGeoService) GetLocation(ip string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	searcher := s.getSearcher(ip)
	if searcher == nil {
		return ""
	}

	region, err := searcher.Search(ip)
	if err != nil {
		log.Printf("[IPGeo] 查询IP失败: %s, error: %v", ip, err)
		return ""
	}

	// 解析结果：国家|区域|省份|城市|ISP
	parts := strings.Split(region, "|")
	if len(parts) < 3 {
		return ""
	}

	country := parts[0]
	province := parts[2]
	city := ""
	if len(parts) >= 4 {
		city = parts[3]
	}

	isChina := country == "中国" || country == "0" || strings.Contains(country, "中国")

	// 海外IP：返回国家名
	if !isChina && country != "0" && country != "" {
		return country
	}

	// 国内IP：确保只返回省级
	// 先规范化（去掉省/市/自治区等后缀）
	provNorm := normalizeRegionName(province)
	cityNorm := normalizeRegionName(city)

	// 情况1：省份字段在白名单中 -> 直接返回（保持原始格式）
	if chinaProvinces[province] {
		return province
	}

	// 情况2：规范化后的省份字段在白名单中 -> 返回规范化后的省份名
	if chinaProvinces[provNorm] {
		return provNorm
	}

	// 情况3：省份字段是地级市名 -> 查映射表得到省份
	if mappedProvince, ok := knownNonProvinceCities[provNorm]; ok {
		return mappedProvince
	}
	// 再试一次原始省份名（兼容映射表里直接写了带后缀的情况）
	if mappedProvince, ok := knownNonProvinceCities[province]; ok {
		return mappedProvince
	}

	// 情况4：省份为空但有城市名 -> 通过城市名反查省份
	if (province == "0" || province == "") && city != "0" && city != "" {
		if mappedProvince, ok := knownNonProvinceCities[cityNorm]; ok {
			return mappedProvince
		}
		if mappedProvince, ok := knownNonProvinceCities[city]; ok {
			return mappedProvince
		}
	}

	// 情况5：省份为空但有国家名（海外兜底）
	if (province == "0" || province == "") && country != "0" && country != "" {
		return country
	}

	// 兜底：返回规范化后的省份名
	if provNorm != "" && provNorm != "0" {
		return provNorm
	}
	return ""
}

// GetCityLocation 获取IP地理位置（城市级别）
// 国内IP返回"省份 城市"，海外IP返回国家名
func (s *IPGeoService) GetCityLocation(ip string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	searcher := s.getSearcher(ip)
	if searcher == nil {
		return ""
	}

	region, err := searcher.Search(ip)
	if err != nil {
		return ""
	}

	parts := strings.Split(region, "|")
	if len(parts) < 3 {
		return ""
	}

	country := parts[0]
	province := parts[2]
	city := ""
	if len(parts) >= 4 {
		city = parts[3]
	}

	isChina := country == "中国" || country == "0" || strings.Contains(country, "中国")

	if !isChina && country != "0" && country != "" {
		return country
	}

	provNorm := normalizeRegionName(province)
	cityNorm := normalizeRegionName(city)

	// 获取有效省份名
	provinceName := ""
	if chinaProvinces[province] {
		provinceName = province
	} else if chinaProvinces[provNorm] {
		provinceName = provNorm
	} else if mapped, ok := knownNonProvinceCities[provNorm]; ok {
		provinceName = mapped
	} else if mapped, ok := knownNonProvinceCities[province]; ok {
		provinceName = mapped
	}

	// 获取有效城市名
	cityName := ""
	if city != "0" && city != "" && cityNorm != provinceName {
		cityName = cityNorm
	}

	if provinceName != "" && cityName != "" {
		return provinceName + " " + cityName
	}
	if provinceName != "" {
		return provinceName
	}
	return ""
}

// GetFullLocation 获取完整IP地理位置
func (s *IPGeoService) GetFullLocation(ip string) map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := map[string]string{
		"country":  "",
		"region":   "",
		"province": "",
		"city":     "",
		"isp":      "",
	}

	searcher := s.getSearcher(ip)
	if searcher == nil {
		return result
	}

	region, err := searcher.Search(ip)
	if err != nil {
		log.Printf("[IPGeo] 查询IP失败: %s, error: %v", ip, err)
		return result
	}

	// 解析结果：国家-区域|省份|城市|ISP
	parts := strings.Split(region, "|")
	if len(parts) >= 5 {
		for i, key := range []string{"country", "region", "province", "city", "isp"} {
			if parts[i] != "0" {
				result[key] = parts[i]
			}
		}
	}

	return result
}

// UpdateDatabase 更新数据库
func (s *IPGeoService) UpdateDatabase() error {
	config := s.loadConfig()

	// 确保目录存在
	dir := filepath.Dir(config.DatabasePathV4)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 下载 IPv4 数据库
	log.Printf("[IPGeo] 开始下载IPv4数据库: %s", config.DownloadURLV4)
	if err := s.downloadFile(config.DownloadURLV4, config.DatabasePathV4); err != nil {
		return fmt.Errorf("下载IPv4数据库失败: %v", err)
	}

	// 下载 IPv6 数据库
	log.Printf("[IPGeo] 开始下载IPv6数据库: %s", config.DownloadURLV6)
	if err := s.downloadFile(config.DownloadURLV6, config.DatabasePathV6); err != nil {
		log.Printf("[IPGeo] 下载IPv6数据库失败（可选）: %v", err)
	}

	// 重新加载数据库
	if err := s.loadDatabase(); err != nil {
		return fmt.Errorf("重新加载数据库失败: %v", err)
	}

	// 更新配置
	config.LastUpdate = time.Now().Format("2006-01-02 15:04:05")
	if err := s.saveConfig(config); err != nil {
		log.Printf("[IPGeo] 保存配置失败: %v", err)
	}

	log.Printf("[IPGeo] 数据库更新成功")
	return nil
}

// downloadFile 下载文件
func (s *IPGeoService) downloadFile(url, filePath string) error {
	// === 下载前校验（SSRF 防护）===
	if err := validateDownloadURL(url); err != nil {
		return fmt.Errorf("下载地址校验失败: %v", err)
	}

	// 创建临时文件
	tmpFile := filePath + ".tmp"
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// 使用自定义超时的 HTTP 客户端
	client := &http.Client{
		Timeout: 300 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		os.Remove(tmpFile)
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	resp, err := client.Do(req)
	if err != nil {
		os.Remove(tmpFile)
		return err
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		os.Remove(tmpFile)
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 限制读取大小（防止 zip bomb / 资源耗尽）
	limitedReader := io.LimitReader(resp.Body, maxXDBFileSize+1024*1024)
	written, err := io.Copy(out, limitedReader)
	out.Close()
	if err != nil {
		os.Remove(tmpFile)
		return err
	}
	if written > maxXDBFileSize {
		os.Remove(tmpFile)
		return fmt.Errorf("文件超过大小限制: %d bytes", written)
	}

	// === 下载后校验（防止下载到恶意文件）===
	if err := validateDownloadedFile(tmpFile); err != nil {
		os.Remove(tmpFile) // 校验失败，删除恶意文件
		return fmt.Errorf("文件校验失败（可能不是合法数据库）: %v", err)
	}

	// 替换原文件
	if err := os.Rename(tmpFile, filePath); err != nil {
		os.Remove(tmpFile)
		return err
	}

	return nil
}

// startAutoUpdate 启动自动更新检查
func (s *IPGeoService) startAutoUpdate() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAndUpdate()
		case <-s.stopChan:
			return
		}
	}
}

// checkAndUpdate 检查并更新
func (s *IPGeoService) checkAndUpdate() {
	config := s.loadConfig()

	if !config.Enabled {
		return
	}

	// 检查是否需要更新
	if config.LastUpdate == "" {
		if err := s.UpdateDatabase(); err != nil {
			log.Printf("[IPGeo] 自动更新失败: %v", err)
		}
		return
	}

	// 解析上次更新时间
	lastUpdate, err := time.Parse("2006-01-02 15:04:05", config.LastUpdate)
	if err != nil {
		log.Printf("[IPGeo] 解析更新时间失败: %v", err)
		return
	}

	// 检查是否达到更新间隔
	days := config.UpdateDays
	if days <= 0 {
		days = 10
	}
	if time.Since(lastUpdate) >= time.Duration(days)*24*time.Hour {
		log.Printf("[IPGeo] 达到更新间隔，开始自动更新")
		if err := s.UpdateDatabase(); err != nil {
			log.Printf("[IPGeo] 自动更新失败: %v", err)
		}
	}
}

// GetConfig 获取配置
func (s *IPGeoService) GetConfig() IPGeoConfig {
	return s.loadConfig()
}

// UpdateConfig 更新配置
func (s *IPGeoService) UpdateConfig(config IPGeoConfig) error {
	// === 前置校验（SSRF 防护）：保存前校验下载地址 ===
	if config.DownloadURLV4 != "" {
		if err := validateDownloadURL(config.DownloadURLV4); err != nil {
			return fmt.Errorf("IPv4下载地址不安全: %v", err)
		}
	}
	if config.DownloadURLV6 != "" {
		if err := validateDownloadURL(config.DownloadURLV6); err != nil {
			return fmt.Errorf("IPv6下载地址不安全: %v", err)
		}
	}

	// 保存配置
	if err := s.saveConfig(config); err != nil {
		return err
	}

	// 如果启用了数据库，尝试加载
	if config.Enabled {
		s.dbPathV4 = config.DatabasePathV4
		s.dbPathV6 = config.DatabasePathV6
		if err := s.loadDatabase(); err != nil {
			log.Printf("[IPGeo] 加载数据库失败: %v", err)
		}
	}

	return nil
}

// Stop 停止服务
func (s *IPGeoService) Stop() {
	close(s.stopChan)
	if s.searcherV4 != nil {
		s.searcherV4.Close()
	}
	if s.searcherV6 != nil {
		s.searcherV6.Close()
	}
}

// IsEnabled 检查是否启用
func (s *IPGeoService) IsEnabled() bool {
	config := s.loadConfig()
	return config.Enabled
}

// GetDatabaseInfo 获取数据库信息
func (s *IPGeoService) GetDatabaseInfo() map[string]interface{} {
	config := s.loadConfig()
	info := map[string]interface{}{
		"enabled":            config.Enabled,
		"download_url_v4":    config.DownloadURLV4,
		"download_url_v6":    config.DownloadURLV6,
		"update_days":        config.UpdateDays,
		"last_update":        config.LastUpdate,
		"database_path_v4":   config.DatabasePathV4,
		"database_path_v6":   config.DatabasePathV6,
		"ipv4_loaded":        s.searcherV4 != nil,
		"ipv6_loaded":        s.searcherV6 != nil,
	}

	// 检查 IPv4 文件
	if stat, err := os.Stat(config.DatabasePathV4); err == nil {
		info["file_size_v4"] = stat.Size()
		info["file_exists_v4"] = true
	} else {
		info["file_exists_v4"] = false
	}

	// 检查 IPv6 文件
	if stat, err := os.Stat(config.DatabasePathV6); err == nil {
		info["file_size_v6"] = stat.Size()
		info["file_exists_v6"] = true
	} else {
		info["file_exists_v6"] = false
	}

	return info
}
