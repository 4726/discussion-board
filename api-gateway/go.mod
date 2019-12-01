module github.com/4726/discussion-board/api-gateway

go 1.13

require (
	github.com/4726/discussion-board/services/common v0.0.0-00010101000000-000000000000
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.4.0
	github.com/golang/protobuf v1.3.2
	github.com/kr/pretty v0.1.0 // indirect
	github.com/stretchr/testify v1.3.0
	golang.org/x/text v0.3.1-0.20180807135948-17ff2d5776d2 // indirect
	google.golang.org/genproto v0.0.0-20190404172233-64821d5d2107 // indirect
	google.golang.org/grpc v1.19.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/4726/discussion-board/services/common => ../services/common
