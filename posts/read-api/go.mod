module github.com/4726/discussion-board/posts/read-api

go 1.13

require (
	github.com/4726/discussion-board/posts/models v0.0.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jinzhu/gorm v1.9.11
	github.com/kr/pretty v0.1.0 // indirect
	github.com/stretchr/testify v1.3.0
)

replace github.com/4726/discussion-board/posts/models => ../models
