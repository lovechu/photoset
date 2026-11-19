package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"photoset/internal/database"
	"photoset/internal/repository"
)

const (
	SiteBaseURL       = "https://tt.cy.mk"
	sitemapCacheKey   = "photoset:sitemap:"
	sitemapCacheTTL   = time.Hour
	maxURLsPerMap     = 50000
	maxImagePerURL    = 1000
)

// ---- XML 结构定义 ----

type urlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	XmlnsImg string  `xml:"xmlns:image,attr,omitempty"`
	URLs    []urlEntry `xml:"url"`
}

type urlEntry struct {
	Loc        string       `xml:"loc"`
	LastMod    string       `xml:"lastmod,omitempty"`
	ChangeFreq string       `xml:"changefreq,omitempty"`
	Priority   string       `xml:"priority,omitempty"`
	Images     []imageEntry `xml:"image:image,omitempty"`
}

type imageEntry struct {
	Loc string `xml:"image:loc"`
}

type sitemapIndex struct {
	XMLName xml.Name      `xml:"sitemapindex"`
	Xmlns   string        `xml:"xmlns,attr"`
	Maps    []sitemapRef  `xml:"sitemap"`
}

type sitemapRef struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

// SitemapService 生成并缓存 sitemap
type SitemapService struct {
	repo *repository.SeoRepository
}

func NewSitemapService() *SitemapService {
	return &SitemapService{repo: repository.NewSeoRepository()}
}

func (s *SitemapService) cached(ctx context.Context, key string, gen func() (string, error)) (string, error) {
	cache := NewCacheService()
	var out string
	if database.IsRedisAvailable() {
		if err := cache.Get(ctx, key, &out); err == nil && out != "" {
			return out, nil
		}
	}
	content, err := gen()
	if err != nil {
		return "", err
	}
	if database.IsRedisAvailable() {
		_ = cache.Set(ctx, key, content, sitemapCacheTTL)
	}
	return content, nil
}

// InvalidateCache 清除全部 sitemap 缓存（发布/删除内容时调用）
func (s *SitemapService) InvalidateCache(ctx context.Context) {
	cache := NewCacheService()
	_ = cache.DeleteByPattern(ctx, sitemapCacheKey+"*")
}

// Robots 返回抓取规则
func (s *SitemapService) Robots() string {
	return fmt.Sprintf(`User-agent: *
Allow: /
Disallow: /admin/
Disallow: /login
Disallow: /register
Disallow: /api/
Disallow: /user/
Disallow: /profile
Disallow: /edit-profile
Disallow: /orders
Disallow: /points
Disallow: /trash
Disallow: /search

Sitemap: %s/sitemap_index.xml
`, SiteBaseURL)
}

// Index 站点地图索引
func (s *SitemapService) Index(ctx context.Context) (string, error) {
	return s.cached(ctx, sitemapCacheKey+"index", func() (string, error) {
		now := time.Now().UTC().Format("2006-01-02")
		idx := sitemapIndex{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
		for _, name := range []string{"static", "photosets", "posts"} {
			idx.Maps = append(idx.Maps, sitemapRef{
				Loc:     fmt.Sprintf("%s/sitemap-%s.xml", SiteBaseURL, name),
				LastMod: now,
			})
		}
		b, err := xml.MarshalIndent(idx, "", "  ")
		return xml.Header + string(b), err
	})
}

// Static 固定页面 sitemap
func (s *SitemapService) Static(ctx context.Context) (string, error) {
	return s.cached(ctx, sitemapCacheKey+"static", func() (string, error) {
		now := time.Now().UTC().Format("2006-01-02")
		paths := []struct {
			p string
			f string
		}{
			{"/", "daily"},
			{"/explore", "daily"},
			{"/community", "daily"},
			{"/community/hot", "hourly"},
			{"/tags", "weekly"},
			{"/leaderboard", "daily"},
		}
		us := urlSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
		for _, item := range paths {
			us.URLs = append(us.URLs, urlEntry{
				Loc:        SiteBaseURL + item.p,
				LastMod:    now,
				ChangeFreq: item.f,
				Priority:   "0.8",
			})
		}
		// 分类页
		cats, _ := s.repo.ListCategories()
		for _, c := range cats {
			us.URLs = append(us.URLs, urlEntry{
				Loc:        fmt.Sprintf("%s/?category=%s", SiteBaseURL, c),
				ChangeFreq: "weekly",
				Priority:   "0.6",
			})
		}
		b, err := xml.MarshalIndent(us, "", "  ")
		return xml.Header + string(b), err
	})
}

// Photosets 套图 sitemap，含图片扩展信息（<image:image> 合并进来）
func (s *SitemapService) Photosets(ctx context.Context) (string, error) {
	return s.cached(ctx, sitemapCacheKey+"photosets", func() (string, error) {
		items, err := s.repo.ListPublishedPhotosets()
		if err != nil {
			return "", err
		}
		us := urlSet{
			Xmlns:    "http://www.sitemaps.org/schemas/sitemap/0.9",
			XmlnsImg: "http://www.google.com/schemas/sitemap-image/1.1",
		}
		count := 0
		for _, sp := range items {
			if count >= maxURLsPerMap {
				break
			}
			entry := urlEntry{
				Loc:        fmt.Sprintf("%s/photoset/%d", SiteBaseURL, sp.ID),
				LastMod:    sp.UpdatedAt.UTC().Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   "0.9",
			}
			for i, img := range sp.ImageURLs {
				if i >= maxImagePerURL {
					break
				}
				loc := img
				if strings.HasPrefix(loc, "/") {
					loc = SiteBaseURL + loc
				} else if !strings.HasPrefix(loc, "http") {
					loc = SiteBaseURL + "/" + loc
				}
				entry.Images = append(entry.Images, imageEntry{Loc: loc})
			}
			us.URLs = append(us.URLs, entry)
			count++
		}
		b, err := xml.MarshalIndent(us, "", "  ")
		return xml.Header + string(b), err
	})
}

// Posts 帖子 sitemap
func (s *SitemapService) Posts(ctx context.Context) (string, error) {
	return s.cached(ctx, sitemapCacheKey+"posts", func() (string, error) {
		items, err := s.repo.ListPublicPosts()
		if err != nil {
			return "", err
		}
		us := urlSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
		count := 0
		for _, p := range items {
			if count >= maxURLsPerMap {
				break
			}
			us.URLs = append(us.URLs, urlEntry{
				Loc:        fmt.Sprintf("%s/community/post/%d", SiteBaseURL, p.ID),
				LastMod:    p.UpdatedAt.UTC().Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   "0.7",
			})
			count++
		}
		b, err := xml.MarshalIndent(us, "", "  ")
		return xml.Header + string(b), err
	})
}
