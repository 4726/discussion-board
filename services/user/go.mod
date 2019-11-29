module github.com/4726/discussion-board/services/user

go 1.13

require (
	github.com/4726/discussion-board/services/common v0.0.0-20191128014125-0a385bcc49c5
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.2
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/jinzhu/gorm v1.9.11
	github.com/kr/pretty v0.1.0 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20191028145041-f83a4685e152
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.25.1
)

replace github.com/4726/discussion-board/services/common => ../common
