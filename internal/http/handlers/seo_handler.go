package handlers

import (
	"net/http"

	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

// SeoHandler robots.txt 与 sitemap 输出
type SeoHandler struct {
	svc *service.SitemapService
}

func NewSeoHandler() *SeoHandler {
	return &SeoHandler{svc: service.NewSitemapService()}
}

func (h *SeoHandler) Robots(c *gin.Context) {
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=3600")
	c.String(http.StatusOK, h.svc.Robots())
}

func (h *SeoHandler) Index(c *gin.Context) {
	h.serveXML(c, func() (string, error) { return h.svc.Index(c.Request.Context()) })
}

func (h *SeoHandler) Static(c *gin.Context) {
	h.serveXML(c, func() (string, error) { return h.svc.Static(c.Request.Context()) })
}

func (h *SeoHandler) Photosets(c *gin.Context) {
	h.serveXML(c, func() (string, error) { return h.svc.Photosets(c.Request.Context()) })
}

func (h *SeoHandler) Posts(c *gin.Context) {
	h.serveXML(c, func() (string, error) { return h.svc.Posts(c.Request.Context()) })
}

func (h *SeoHandler) serveXML(c *gin.Context, gen func() (string, error)) {
	content, err := gen()
	if err != nil {
		c.String(http.StatusServiceUnavailable, "")
		return
	}
	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=3600")
	c.String(http.StatusOK, content)
}
