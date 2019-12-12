module github.com/4726/discussion-board/services/media

go 1.13

require (
	github.com/4726/discussion-board/services/common v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/jinzhu/gorm v1.9.11
	github.com/minio/minio-go/v6 v6.0.39
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/segmentio/ksuid v1.0.2
	github.com/stretchr/testify v1.3.0
	google.golang.org/grpc v1.25.1
)

replace github.com/4726/discussion-board/services/common => ../common
