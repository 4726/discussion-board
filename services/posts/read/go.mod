module github.com/4726/discussion-board/posts/services/read

go 1.13

require (
	github.com/4726/discussion-board/services/common v0.0.0-00010101000000-000000000000
	github.com/4726/discussion-board/services/posts/models v0.0.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jinzhu/gorm v1.9.11
	github.com/kr/pretty v0.1.0 // indirect
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.3.0
)

replace github.com/4726/discussion-board/services/posts/models => ../models

replace github.com/4726/discussion-board/services/common => ../../common
