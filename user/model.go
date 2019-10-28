package main

import (
	"time"
)

type Auth struct {
	UserID    int    `gorm:"AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	Username  string `gorm:"UNIQUE_INDEX;type:varchar(128)"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Profile struct {
	UserID   int    `gorm:"AUTO_INCREMENT:false;NOT NULL;PRIMARY_KEY"`
	Username string `gorm:"UNIQUE_INDEX;type:varchar(128)"`
	Bio      string
	AvatarID string
}
