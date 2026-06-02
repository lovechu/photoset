package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"photoset/internal/pkg/jwt"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// wsAllowedOrigins 从环境变量读取 WebSocket 允许的来源
func wsAllowedOrigins() []string {
	origins := os.Getenv("CORS_ALLOW_ORIGINS")
	if origins == "" {
		return []string{}
	}
	var result []string
	for _, o := range strings.Split(origins, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			result = append(result, o)
		}
	}
	return result
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 开发环境（未配置 CORS_ALLOW_ORIGINS）允许所有来源
		allowed := wsAllowedOrigins()
		if len(allowed) == 0 {
			return true
		}

		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}

		for _, o := range allowed {
			if o == origin {
				return true
			}
		}

		log.Printf("[WebSocket] Origin %s not allowed", origin)
		return false
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub *service.Hub
}

// NewWebSocketHandler creates a new WebSocketHandler
func NewWebSocketHandler(hub *service.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleConnection upgrades HTTP to WebSocket and manages the connection
// GET /api/community/ws
// 认证方式: URL query 参数 ?token=xxx 或 Authorization header
// 升级后客户端也可发送 auth 消息进行认证
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	// 从 query 参数或 Authorization header 提取 token
	token := extractWSToken(c)

	var userID uint
	if token != "" {
		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "token 无效或已过期"})
			return
		}
		userID = claims.UserID
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WebSocket] Upgrade error: %v", err)
		return
	}

	// 如果连接时没有 token，等待客户端发送 auth 消息
	if userID == 0 {
		userID = h.waitForAuth(conn)
		if userID == 0 {
			// 认证失败，关闭连接
			conn.WriteJSON(service.WSMessage{
				Type:    "error",
				Payload: map[string]string{"message": "认证失败，请发送有效的 auth 消息"},
			})
			conn.Close()
			return
		}
	}

	client := h.hub.NewClient(conn, userID)
	h.hub.Register(client)

	// 发送认证成功消息
	client.SendJSON(service.WSMessage{
		Type:    "auth_ok",
		Payload: map[string]interface{}{"user_id": userID},
	})

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}

// extractWSToken 从 HTTP 请求中提取 JWT token
// 优先从 Authorization header 提取，其次从 query 参数提取
func extractWSToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}
	return c.Query("token")
}

// waitForAuth 等待客户端发送 auth 消息，返回 userID（0 表示失败）
func (h *WebSocketHandler) waitForAuth(conn *websocket.Conn) uint {
	// 设置认证超时 10 秒
	conn.SetReadLimit(4096)
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Printf("[WebSocket] Read auth message error: %v", err)
		return 0
	}

	var msg service.WSMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("[WebSocket] Parse auth message error: %v", err)
		return 0
	}

	if msg.Type != "auth" {
		log.Printf("[WebSocket] Expected auth message, got: %s", msg.Type)
		return 0
	}

	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		log.Printf("[WebSocket] Invalid auth payload format")
		return 0
	}

	tokenRaw, ok := payload["token"]
	if !ok {
		log.Printf("[WebSocket] Missing token in auth payload")
		return 0
	}

	token, ok := tokenRaw.(string)
	if !ok || token == "" {
		log.Printf("[WebSocket] Token is not a string or empty")
		return 0
	}

	claims, err := jwt.ParseToken(token)
	if err != nil {
		log.Printf("[WebSocket] Invalid token in auth message: %v", err)
		return 0
	}

	return claims.UserID
}

// GetOnlineStatus returns online status of specific users
// GET /api/community/ws/online?user_ids=1,2,3
func (h *WebSocketHandler) GetOnlineStatus(c *gin.Context) {
	userIDsStr := c.QueryArray("user_ids")
	if len(userIDsStr) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "user_ids is required"})
		return
	}

	type OnlineStatus struct {
		UserID uint `json:"user_id"`
		Online bool `json:"online"`
	}

	var statuses []OnlineStatus
	for _, idStr := range userIDsStr {
		var id uint
		_, _ = fmt.Sscanf(idStr, "%d", &id)
		if id > 0 {
			statuses = append(statuses, OnlineStatus{
				UserID: id,
				Online: h.hub.IsUserOnline(id),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"statuses": statuses,
			"total_online": h.hub.TotalConnections(),
		},
	})
}
