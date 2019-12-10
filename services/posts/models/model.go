package models

import (
	"time"
)

type Comment struct {
	ID        uint64 `gorm:"AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	PostID    uint64
	ParentID  uint64 `gorm:"DEFAULT:0"`
	UserID    uint64
	Body      string
	CreatedAt time.Time
	Likes     int64
}

type Post struct {
	ID        uint64 `gorm:"AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	UserID    uint64
	Title     string
	Body      string
	Likes     int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Comments  []Comment
}
