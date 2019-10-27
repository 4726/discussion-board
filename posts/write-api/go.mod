module github.com/4726/discussion-board/posts/write-api

go 1.13

require (
	github.com/4726/discussion-board/posts/models v0.0.0
	github.com/gin-gonic/gin v1.4.0 // indirect
	github.com/jinzhu/gorm v1.9.11 // indirect
)

replace github.com/4726/discussion-board/posts/models => ../models
