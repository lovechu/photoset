package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"photoset/internal/database"

	"github.com/redis/go-redis/v9"
)

const seoRenderStream = "photoset:seo:render"

// SeoEventService 发布 SEO 事件到 Redis Stream（持久化，避免 fire-and-forget 丢消息）。
// renderer 以消费者组消费本 Stream，保证内容更新后预渲染/失效事件可靠送达。
type SeoEventService struct {
	sitemap *SitemapService
}

func NewSeoEventService() *SeoEventService {
	return &SeoEventService{sitemap: NewSitemapService()}
}

// publish 向 Stream 写入一条事件。Redis 不可用时仅告警，绝不影响主业务流程。
func (s *SeoEventService) publish(action, path string) {
	if !database.IsRedisAvailable() || database.RedisClient == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	values := []interface{}{
		"action", action, // render | invalidate
		"path", path,
		"ts", time.Now().Unix(),
	}
	if err := database.RedisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: seoRenderStream,
		Values: values,
	}).Err(); err != nil {
		log.Printf("[SEO] 发布渲染事件失败 action=%s path=%s: %v", action, path, err)
	}
}

// InvalidateSitemap 清除 sitemap 缓存（发布/更新/删除后调用）
func (s *SeoEventService) InvalidateSitemap() {
	s.sitemap.InvalidateCache(context.Background())
}

// PhotosetChanged 套图创建/更新：失效 sitemap + 请求预渲染
func (s *SeoEventService) PhotosetChanged(id uint) {
	s.InvalidateSitemap()
	s.publish("render", fmt.Sprintf("/photoset/%d", id))
}

// PhotosetRemoved 套图删除/永久删除：失效 sitemap + 请求清缓存
func (s *SeoEventService) PhotosetRemoved(id uint) {
	s.InvalidateSitemap()
	s.publish("invalidate", fmt.Sprintf("/photoset/%d", id))
}

// PostChanged 帖子创建/更新：失效 sitemap + 请求预渲染
func (s *SeoEventService) PostChanged(id uint) {
	s.InvalidateSitemap()
	s.publish("render", fmt.Sprintf("/community/post/%d", id))
}

// PostRemoved 帖子删除：失效 sitemap + 请求清缓存
func (s *SeoEventService) PostRemoved(id uint) {
	s.InvalidateSitemap()
	s.publish("invalidate", fmt.Sprintf("/community/post/%d", id))
}
