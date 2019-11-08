package models

import (
	"time"
)

type Comment struct {
	ID        uint `gorm:"AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	PostID    uint
	ParentID  uint `gorm:"DEFAULT:0"`
	User      string
	Body      string
	CreatedAt time.Time
	Likes     int
}

type Post struct {
	ID        uint `gorm:"AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	User      string
	Title     string
	Body      string
	Likes     int
	CreatedAt time.Time
	UpdatedAt time.Time
	Comments  []Comment
}