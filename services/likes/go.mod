module github.com/4726/discussion-board/services/likes

go 1.13

require (
	github.com/4726/discussion-board/services/common v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jinzhu/gorm v1.9.11
	github.com/kr/pretty v0.1.0 // indirect
	github.com/prometheus/client_golang v1.2.1
	github.com/stretchr/testify v1.4.0
)

replace github.com/4726/discussion-board/services/common => ../common
