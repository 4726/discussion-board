module github.com/4726/discussion-board/services/media

go 1.13

require (
	github.com/4726/discussion-board/services/common v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.2
	github.com/jinzhu/gorm v1.9.11
	github.com/kr/pretty v0.1.0 // indirect
	github.com/minio/minio-go/v6 v6.0.39
	github.com/prometheus/client_golang v1.2.1 // indirect
	github.com/segmentio/ksuid v1.0.2
	github.com/stretchr/testify v1.3.0
	google.golang.org/grpc v1.25.1
)

replace github.com/4726/discussion-board/services/common => ../common
