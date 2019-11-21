module github.com/4726/discussion-board/api-gateway

go 1.13

require (
	github.com/4726/discussion-board/services/common v0.0.0-00010101000000-000000000000
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/jinzhu/gorm v1.9.11
	github.com/stretchr/testify v1.3.0
)

replace github.com/4726/discussion-board/services/common => ../services/common
