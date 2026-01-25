package handler

import (
	"database/sql"
	"net/http"
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
	db *sql.DB
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(db *sql.DB) *NotificationHandler {
	h := &NotificationHandler{db: db}
	return h
}

// ensureDefaultNotifications creates default notifications if none exist for the user
func (h *NotificationHandler) ensureDefaultNotifications(userID string) {
	var count int
	h.db.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = ?", userID).Scan(&count)

	if count == 0 {
		// Add demo notifications for new users
		now := time.Now()
		notifications := []struct {
			Type      string
			Title     string
			Message   string
			ActionURL string
		}{
			{"low_stock", "Low Stock Alert", "Kopi Susu Gula Aren stock is low (10 remaining)", "/products"},
			{"system", "Welcome to GoPOS!", "Your POS system is ready to use.", ""},
		}

		for _, n := range notifications {
			h.db.Exec(`
				INSERT INTO notifications (id, user_id, type, title, message, is_read, action_url, created_at)
				VALUES (?, ?, ?, ?, ?, 0, ?, ?)
			`, uuid.New().String(), userID, n.Type, n.Title, n.Message, n.ActionURL, now)
		}
	}
}

// GetNotifications handles GET /api/v1/notifications
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	// Ensure user has default notifications
	h.ensureDefaultNotifications(claims.UserID)

	unreadOnly := c.Query("unread_only") == "true"

	query := `
		SELECT id, type, title, message, is_read, COALESCE(action_url, ''), created_at
		FROM notifications 
		WHERE user_id = ?
	`
	if unreadOnly {
		query += " AND is_read = 0"
	}
	query += " ORDER BY created_at DESC LIMIT 50"

	rows, err := h.db.QueryContext(c.Request.Context(), query, claims.UserID)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	defer rows.Close()

	notifications := make([]*Notification, 0)
	unreadCount := 0

	for rows.Next() {
		n := &Notification{}
		var createdAt time.Time
		var isRead int
		if err := rows.Scan(&n.ID, &n.Type, &n.Title, &n.Message, &isRead, &n.ActionURL, &createdAt); err != nil {
			continue
		}
		n.IsRead = isRead == 1
		n.CreatedAt = createdAt.Format(time.RFC3339)
		if !n.IsRead {
			unreadCount++
		}
		notifications = append(notifications, n)
	}

	utils.SuccessResponse(c, http.StatusOK, "Notifications retrieved", gin.H{
		"unread_count":  unreadCount,
		"notifications": notifications,
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

	_, err := h.db.ExecContext(c.Request.Context(),
		"UPDATE notifications SET is_read = 1 WHERE id = ? AND user_id = ?",
		notifID, claims.UserID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
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

	_, err := h.db.ExecContext(c.Request.Context(),
		"UPDATE notifications SET is_read = 1 WHERE user_id = ?",
		claims.UserID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "All notifications marked as read", nil)
}

// DeleteNotification handles DELETE /api/v1/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	notifID := c.Param("id")
	claims := middleware.GetCurrentUser(c)
	if claims == nil {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	_, err := h.db.ExecContext(c.Request.Context(),
		"DELETE FROM notifications WHERE id = ? AND user_id = ?",
		notifID, claims.UserID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notification deleted", nil)
}
