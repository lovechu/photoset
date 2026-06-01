package handlers

import (
	"fmt"
	"log"
	"net/http"

	"photoset/internal/http/middleware"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development; restrict in production
	CheckOrigin: func(r *http.Request) bool {
		return true
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
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "请先登录"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WebSocket] Upgrade error: %v", err)
		return
	}

	client := h.hub.NewClient(conn, userID)
	h.hub.Register(client)

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
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
