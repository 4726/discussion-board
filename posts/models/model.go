package models

import (
	"time"
)

type Comment struct {
	CommentID int `gorm:"AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	PostID    int `gorm:"FOREIGNKEY:PostID"`
	ParentID  int `gorm:"DEFAULT:0"`
	User      string
	Body      string
	CreatedAt time.Time
	Likes     int
}

type Post struct {
	PostID    int `gorm:"AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	User      string
	Title     string
	Body      string
	Likes     int
	CreatedAt time.Time
	UpdatedAt time.Time
	Comments  []Comment `gorm:"FOREIGNKEY:PostID`
}
