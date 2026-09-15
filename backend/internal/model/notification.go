package model

import "time"

// Notification merepresentasikan satu baris di tabel "notifications".
type Notification struct {
	ID      int64     `json:"id" gorm:"column:id;primaryKey"`
	UserID  int64     `json:"user_id" gorm:"column:user_id"`
	Type    string    `json:"type" gorm:"column:type"`
	Title   string    `json:"title" gorm:"column:title"`
	Message string    `json:"message" gorm:"column:message"`
	IsRead  bool      `json:"is_read" gorm:"column:is_read"`
	SentAt  time.Time `json:"sent_at" gorm:"column:sent_at"`
}

func (Notification) TableName() string {
	return "notifications"
}
