module github.com/4726/discussion-board/services/likes

go 1.13

require (
	github.com/4726/discussion-board/services/common v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.2
	github.com/jinzhu/gorm v1.9.11
	github.com/kr/pretty v0.1.0 // indirect
	github.com/prometheus/client_golang v1.2.1 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20191119073136-fc4aabc6c914 // indirect
	golang.org/x/sys v0.0.0-20191120155948-bd437916bb0e // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20191115221424-83cc0476cb11 // indirect
	google.golang.org/grpc v1.25.1
)

replace github.com/4726/discussion-board/services/common => ../common
