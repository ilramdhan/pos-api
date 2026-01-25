package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ilramdhan/pos-api/internal/middleware"
	"github.com/ilramdhan/pos-api/internal/utils"
)

// Notification represents a notification
type Notification struct {
	ID        string `json:"id"`
	UserID    string `json:"-"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
	ActionURL string `json:"action_url,omitempty"`
}

// NotificationHandler handles notification endpoints
type NotificationHandler struct {
	mu            sync.RWMutex
	notifications map[string][]*Notification
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler() *NotificationHandler {
	h := &NotificationHandler{
		notifications: make(map[string][]*Notification),
	}
	// Add some default notifications for demo
	h.addDemoNotifications()
	return h
}

func (h *NotificationHandler) addDemoNotifications() {
	demoNotifs := []*Notification{
		{
			ID:        uuid.New().String(),
			Type:      "low_stock",
			Title:     "Low Stock Alert",
			Message:   "Kopi Susu Gula Aren stock is low (10 remaining)",
			IsRead:    false,
			CreatedAt: time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
			ActionURL: "/products",
		},
		{
			ID:        uuid.New().String(),
			Type:      "system",
			Title:     "Welcome to GoPOS!",
			Message:   "Your POS system is ready to use.",
			IsRead:    true,
			CreatedAt: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		},
	}
	h.notifications["demo"] = demoNotifs
}

// GetNotifications handles GET /api/v1/notifications
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	// Get user notifications or demo notifications
	notifs, ok := h.notifications[claims.UserID]
	if !ok {
		notifs = h.notifications["demo"]
	}

	unreadOnly := c.Query("unread_only") == "true"
	var filtered []*Notification
	unreadCount := 0

	for _, n := range notifs {
		if !n.IsRead {
			unreadCount++
		}
		if !unreadOnly || !n.IsRead {
			filtered = append(filtered, n)
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Notifications retrieved", gin.H{
		"unread_count":  unreadCount,
		"notifications": filtered,
	})
}

// MarkAsRead handles PUT /api/v1/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notifID := c.Param("id")
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	notifs := h.notifications[claims.UserID]
	if notifs == nil {
		notifs = h.notifications["demo"]
	}

	for _, n := range notifs {
		if n.ID == notifID {
			n.IsRead = true
			break
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Notification marked as read", nil)
}

// MarkAllAsRead handles PUT /api/v1/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	notifs := h.notifications[claims.UserID]
	if notifs == nil {
		notifs = h.notifications["demo"]
	}

	for _, n := range notifs {
		n.IsRead = true
	}

	utils.SuccessResponse(c, http.StatusOK, "All notifications marked as read", nil)
}
