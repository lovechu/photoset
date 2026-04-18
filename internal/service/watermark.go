package service

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"photoset/internal/repository"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// WatermarkService 水印服务
type WatermarkService struct {
	settingRepo *repository.SiteSettingRepository
	fontData    []byte // 加载的中文字体数据
}

// NewWatermarkService 创建水印服务
func NewWatermarkService() *WatermarkService {
	return &WatermarkService{
		settingRepo: repository.NewSiteSettingRepository(),
		fontData:    getDefaultFont(),
	}
}

// WatermarkConfig 水印配置
type WatermarkConfig struct {
	Enabled  bool
	Text     string
	Opacity  float64 // 0-100
	Position string  // bottom-right, bottom-left, top-right, top-left, center
}

// getWatermarkConfig 获取水印配置
func (s *WatermarkService) getWatermarkConfig() (*WatermarkConfig, error) {
	settings, err := s.settingRepo.GetAll()
	if err != nil {
		return nil, err
	}

	enabled := settings["watermark_enabled"] == "true"
	text := settings["watermark_text"]
	if text == "" {
		text = "© PhotoSet"
	}

	opacity := 30.0
	if opacityStr, ok := settings["watermark_opacity"]; ok {
		if op, err := strconv.ParseFloat(opacityStr, 64); err == nil {
			opacity = math.Min(100, math.Max(0, op))
		}
	}

	position := "bottom-right"
	if pos, ok := settings["watermark_position"]; ok && pos != "" {
		position = pos
	}

	return &WatermarkConfig{
		Enabled:  enabled,
		Text:     text,
		Opacity:  opacity,
		Position: position,
	}, nil
}

// AddWatermark 给图片添加水印
// 返回处理后的图片数据
func (s *WatermarkService) AddWatermark(imgData []byte, contentType string) ([]byte, error) {
	// 检查是否启用水印
	cfg, err := s.getWatermarkConfig()
	if err != nil || !cfg.Enabled || cfg.Text == "" {
		return imgData, nil // 不启用水印或出错时返回原图
	}

	// 解码图片
	img, err := s.decodeImage(imgData, contentType)
	if err != nil {
		return imgData, nil // 解码失败返回原图
	}

	// 创建水印图片
	watermark := s.createWatermarkImage(img.Bounds(), cfg)
	if watermark == nil {
		return imgData, nil
	}

	// 合并水印
	result := s.overlayWatermark(img, watermark, cfg.Position)

	// 编码返回
	return s.encodeImage(result, contentType)
}

// decodeImage 根据 content-type 解码图片
func (s *WatermarkService) decodeImage(data []byte, contentType string) (image.Image, error) {
	reader := bytes.NewReader(data)

	switch {
	case strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg"):
		return jpeg.Decode(reader)
	case strings.Contains(contentType, "png"):
		return png.Decode(reader)
	case strings.Contains(contentType, "webp"):
		return imaging.Decode(reader)
	case strings.Contains(contentType, "gif"):
		return imaging.Decode(reader)
	default:
		// 尝试自动检测
		return imaging.Decode(reader)
	}
}

// encodeImage 根据 content-type 编码图片
func (s *WatermarkService) encodeImage(img image.Image, contentType string) ([]byte, error) {
	var buf bytes.Buffer

	switch {
	case strings.Contains(contentType, "jpeg") || strings.Contains(contentType, "jpg"):
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
		return buf.Bytes(), err
	case strings.Contains(contentType, "png"):
		err := png.Encode(&buf, img)
		return buf.Bytes(), err
	default:
		// 默认用 PNG
		err := png.Encode(&buf, img)
		return buf.Bytes(), err
	}
}

// createWatermarkImage 创建水印图片
func (s *WatermarkService) createWatermarkImage(bounds image.Rectangle, cfg *WatermarkConfig) image.Image {
	// 计算水印文字大小（基于图片宽度的 5%）
	fontSize := float64(bounds.Dx()) * 0.04
	if fontSize < 12 {
		fontSize = 12
	}
	if fontSize > 72 {
		fontSize = 72
	}

	// 确保字体已加载
	if len(s.fontData) == 0 {
		s.fontData = getDefaultFont()
	}
	if len(s.fontData) == 0 {
		return s.createSimpleWatermark(bounds, cfg)
	}

	// 解析字体
	f, err := truetype.Parse(s.fontData)
	if err != nil {
		return s.createSimpleWatermark(bounds, cfg)
	}

	// 创建文字图片
	dpi := 72.0
	drawer := &font.Drawer{
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    fontSize,
			DPI:     dpi,
			Hinting: font.HintingFull,
		}),
	}

	// 计算文字尺寸
	textWidth := drawer.MeasureString(cfg.Text)
	textHeight := fixed.I(int(fontSize))

	// 添加 padding
	padding := int(fontSize * 0.5)
	imgWidth := int(textWidth.Ceil()) + padding*2
	imgHeight := int(textHeight.Ceil()) + padding*2

	// 创建 RGBA 图片（支持透明度）
	watermark := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// 设置透明度背景（使用透明 Uniform）
	alpha := uint8(cfg.Opacity * 255 / 100)
	bgColor := image.NewUniform(color.RGBA{R: 0, G: 0, B: 0, A: 0}) // 透明背景
	draw.Draw(watermark, watermark.Bounds(), bgColor, image.ZP, draw.Src)

	// 绘制文字
	drawer.Dst = watermark
	drawer.Src = image.NewUniform(color.RGBA{R: 255, G: 255, B: 255, A: alpha})
	drawer.Dot = fixed.Point26_6{
		X: fixed.I(padding),
		Y: fixed.I(int(fontSize) + padding/2),
	}
	drawer.DrawString(cfg.Text)

	return watermark
}

// overlayWatermark 将水印叠加到原图
func (s *WatermarkService) overlayWatermark(src image.Image, watermark image.Image, position string) image.Image {
	srcBounds := src.Bounds()
	watermarkBounds := watermark.Bounds()

	// 创建输出图片
	result := image.NewRGBA(srcBounds)

	// 先绘制原图
	draw.Draw(result, srcBounds, src, image.ZP, draw.Src)

	// 计算水印位置
	var offset image.Point
	padding := srcBounds.Dx() / 50 // 边距为图片宽度的 2%
	if padding < 20 {
		padding = 20
	}

	watermarkW := watermarkBounds.Dx()
	watermarkH := watermarkBounds.Dy()

	switch position {
	case "bottom-left":
		offset = image.Point{X: padding, Y: srcBounds.Dy() - watermarkH - padding}
	case "top-right":
		offset = image.Point{X: srcBounds.Dx() - watermarkW - padding, Y: padding}
	case "top-left":
		offset = image.Point{X: padding, Y: padding}
	case "center":
		offset = image.Point{
			X: (srcBounds.Dx() - watermarkW) / 2,
			Y: (srcBounds.Dy() - watermarkH) / 2,
		}
	default: // "bottom-right"
		offset = image.Point{
			X: srcBounds.Dx() - watermarkW - padding,
			Y: srcBounds.Dy() - watermarkH - padding,
		}
	}

	// 叠加水印（使用 Over 模式保留透明度）
	draw.Draw(result, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

	return result
}

// createSimpleWatermark 创建简单的水印图片（无字体时使用纯色方块+文字图案）
func (s *WatermarkService) createSimpleWatermark(bounds image.Rectangle, cfg *WatermarkConfig) image.Image {
	// 简单的水印：使用文字图案（方块 + 简单像素字符）
	padding := 8
	fontSize := float64(bounds.Dx()) * 0.03
	if fontSize < 10 {
		fontSize = 10
	}
	textLen := len(cfg.Text)
	imgWidth := int(float64(textLen)*fontSize*0.6) + padding*2
	imgHeight := int(fontSize) + padding*2

	watermark := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// 透明背景
	bgColor := image.NewUniform(color.RGBA{R: 0, G: 0, B: 0, A: 0})
	draw.Draw(watermark, watermark.Bounds(), bgColor, image.ZP, draw.Src)

	// 白色文字方块
	alpha := uint8(cfg.Opacity * 255 / 100)
	drawColor := color.RGBA{R: 255, G: 255, B: 255, A: alpha}

	// 绘制简化的文字图案（用矩形块代替）
	// 每个字符用 3x5 的像素块近似
	charWidth := int(fontSize * 0.5)
	charHeight := int(fontSize * 0.8)
	for i, ch := range cfg.Text {
		x := padding + i*charWidth
		y := padding

		// 根据字符类型绘制不同图案
		switch {
		case ch >= 'a' && ch <= 'z':
			// 小写字母：画矩形
			drawRect(watermark, x, y, x+charWidth-1, y+charHeight-1, drawColor)
		case ch >= 'A' && ch <= 'Z':
			// 大写字母：画矩形
			drawRect(watermark, x, y, x+charWidth-1, y+charHeight-1, drawColor)
		case ch >= '0' && ch <= '9':
			// 数字：画矩形
			drawRect(watermark, x, y, x+charWidth-1, y+charHeight-1, drawColor)
		default:
			// 其他字符：画矩形
			drawRect(watermark, x, y, x+charWidth-1, y+charHeight-1, drawColor)
		}
	}

	return watermark
}

// drawRect 绘制矩形
func drawRect(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if x >= 0 && y >= 0 && x < img.Bounds().Dx() && y < img.Bounds().Dy() {
				img.Set(x, y, c)
			}
		}
	}
}

// GetWatermarkInfo 获取水印信息（用于调试）
func (s *WatermarkService) GetWatermarkInfo() map[string]interface{} {
	cfg, err := s.getWatermarkConfig()
	if err != nil {
		return map[string]interface{}{
			"enabled": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"enabled":  cfg.Enabled,
		"text":     cfg.Text,
		"opacity":  cfg.Opacity,
		"position": cfg.Position,
		"hasFont":  len(s.fontData) > 0 || getDefaultFont() != nil,
	}
}

// GetFontPath 获取可用的字体路径
func GetFontPath() string {
	fontPaths := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/mnt/c/Windows/Fonts/msyh.ttc",
		"/mnt/c/Windows/Fonts/simsun.ttc",
		"/mnt/c/Windows/Fonts/arial.ttf",
		"/mnt/c/Windows/Fonts/DejaVuSans.ttf",
	}

	for _, path := range fontPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// InitWatermarkService 初始化水印服务并设置字体
func InitWatermarkService() *WatermarkService {
	ws := NewWatermarkService()

	// 尝试加载字体
	if fontPath := GetFontPath(); fontPath != "" {
		if data, err := os.ReadFile(fontPath); err == nil {
			if _, err := truetype.Parse(data); err == nil {
				ws.fontData = data
				fontCache = data
			}
		}
	}

	return ws
}

// AddWatermarkToReader 从 io.Reader 添加水印
func (s *WatermarkService) AddWatermarkToReader(r io.Reader, contentType string) ([]byte, string, error) {
	// 读取所有数据
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, "", fmt.Errorf("读取数据失败: %w", err)
	}

	// 添加水印
	result, err := s.AddWatermark(data, contentType)
	if err != nil {
		return data, contentType, nil // 出错时返回原数据
	}

	return result, contentType, nil
}

// fontCache 字体缓存
var fontCache []byte
var fontCacheErr error

// getDefaultFont 获取默认字体（按顺序尝试多个系统路径）
func getDefaultFont() []byte {
	if fontCache != nil || fontCacheErr != nil {
		return fontCache
	}

	// 尝试多个常见字体路径
	fontPaths := []string{
		// Linux
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/TTF/DejaVuSans.ttf",
		"/System/Library/Fonts/Hiragino Sans GB.ttc", // macOS
		// Windows WSL
		"/mnt/c/Windows/Fonts/msyh.ttc",   // 微软雅黑
		"/mnt/c/Windows/Fonts/simsun.ttc", // 宋体
		"/mnt/c/Windows/Fonts/arial.ttf",
		"/mnt/c/Windows/Fonts/DejaVuSans.ttf",
		// 当前目录
		"./fonts/DejaVuSans.ttf",
		"./fonts/simhei.ttf",
	}

	for _, path := range fontPaths {
		if data, err := os.ReadFile(path); err == nil {
			if _, err := truetype.Parse(data); err == nil {
				fontCache = data
				fmt.Printf("✓ 水印字体加载成功: %s\n", path)
				return data
			}
		}
	}

	fontCacheErr = fmt.Errorf("未找到可用字体")
	return nil
}

// loadFontFromPath 从指定路径加载字体
func loadFontFromPath(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取字体文件失败: %w", err)
	}
	if _, err := truetype.Parse(data); err != nil {
		return nil, fmt.Errorf("字体文件格式错误: %w", err)
	}
	return data, nil
}

// LoadAndCacheFont 加载并缓存字体
func (s *WatermarkService) LoadAndCacheFont(path string) error {
	data, err := loadFontFromPath(path)
	if err != nil {
		return err
	}
	s.fontData = data
	fontCache = data
	return nil
}

// IsEnabled 检查水印是否启用
func (s *WatermarkService) IsEnabled() bool {
	cfg, err := s.getWatermarkConfig()
	if err != nil {
		return false
	}
	return cfg.Enabled
}
