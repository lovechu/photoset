package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"photoset/internal/repository"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

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
			log.Printf("[IPGeo] 加载数据库失败: %v", err)
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

// GetLocation 获取IP地理位置（省份）
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

	// 解析结果：国家-区域|省份|城市|ISP
	parts := strings.Split(region, "|")
	if len(parts) >= 3 {
		province := parts[2]
		if province == "0" {
			if parts[0] != "0" {
				return parts[0]
			}
			return ""
		}
		return province
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
	// 创建临时文件
	tmpFile := filePath + ".tmp"
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// 使用自定义超时的 HTTP 客户端
	client := &http.Client{
		Timeout: 120 * time.Second,
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

	// 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(tmpFile)
		return err
	}

	// 关闭文件
	out.Close()

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
