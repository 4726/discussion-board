package main

import (
	"time"
)

type PostLike struct {
	PostID    uint `gorm:"primary_key auto_increment:false"`
	UserID    uint `gorm:"primary_key auto_increment:false"`
	CreatedAt time.Time
}

type CommentLike struct {
	CommentID uint `gorm:"primary_key auto_increment:false"`
	UserID    uint `gorm:"primary_key auto_increment:false"`
	CreatedAt time.Time
}
