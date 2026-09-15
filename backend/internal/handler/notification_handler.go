package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"tourism-booking/internal/middleware"
	"tourism-booking/internal/service"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(s service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: s}
}

// List mengembalikan daftar notifikasi milik user yang sedang login beserta jumlah unread.
func (h *NotificationHandler) List(c *gin.Context) {
	val, exists := c.Get(middleware.ContextUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "belum login"})
		return
	}
	userID := val.(int64)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifications, unread, err := h.notificationService.GetNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "gagal mengambil notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"notifications": notifications,
		"unread_count":  unread,
	})
}

// MarkAsRead menandai satu notifikasi sebagai telah dibaca.
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	val, exists := c.Get(middleware.ContextUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "belum login"})
		return
	}
	userID := val.(int64)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ID notifikasi tidak valid"})
		return
	}

	if err := h.notificationService.MarkAsRead(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "gagal memperbarui status notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "notifikasi berhasil ditandai dibaca"})
}

// MarkAllAsRead menandai seluruh notifikasi user sebagai telah dibaca.
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	val, exists := c.Get(middleware.ContextUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "belum login"})
		return
	}
	userID := val.(int64)

	if err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "gagal memperbarui status notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "semua notifikasi berhasil ditandai dibaca"})
}

// Stream melayani Server-Sent Events (SSE) streaming notifikasi real-time.
func (h *NotificationHandler) Stream(c *gin.Context) {
	val, exists := c.Get(middleware.ContextUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "belum login"})
		return
	}
	userID := val.(int64)

	// Header khusus SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Flush()

	ch, unsubscribe := h.notificationService.Subscribe(userID)
	defer unsubscribe()

	// Kirim initial connected event
	c.Stream(func(w io.Writer) bool {
		c.SSEvent("connected", gin.H{"status": "connected", "user_id": userID, "time": time.Now()})
		return false
	})

	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			return false
		case notif, ok := <-ch:
			if !ok {
				return false
			}
			dataBytes, err := json.Marshal(notif)
			if err != nil {
				return true
			}
			fmt.Fprintf(w, "event: notification\ndata: %s\n\n", string(dataBytes))
			c.Writer.Flush()
			return true
		case <-ticker.C:
			fmt.Fprintf(w, "event: ping\ndata: {\"keepalive\": true}\n\n")
			c.Writer.Flush()
			return true
		}
	})
}
