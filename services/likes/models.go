package main

import (
	"time"
)

type PostLike struct {
	PostID    uint64 `gorm:"primary_key auto_increment:false"`
	UserID    uint64 `gorm:"primary_key auto_increment:false"`
	CreatedAt time.Time
}

type CommentLike struct {
	CommentID uint64 `gorm:"primary_key auto_increment:false"`
	UserID    uint64 `gorm:"primary_key auto_increment:false"`
	CreatedAt time.Time
}
