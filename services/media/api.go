package main

import (
	"fmt"
	"net"

	pb "github.com/4726/discussion-board/services/media/pb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/minio/minio-go/v6"
	"google.golang.org/grpc"
	otgrpc "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
)

const (
	bucketExistsErrMsg = "Your previous request to create the named bucket succeeded and you already own it."
)

var bucketName string

type Api struct {
	mc   *minio.Client
	grpc *grpc.Server
}

func NewApi(cfg Config) (*Api, error) {
	mc, err := initMinio(cfg)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.UnaryServerInterceptor()))
	handlers := &Handlers{mc}
	pb.RegisterMediaServer(server, handlers)

	return &Api{mc, server}, nil
}

func (a *Api) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return a.grpc.Serve(lis)
}

func initMinio(cfg Config) (*minio.Client, error) {
	bucketName = cfg.BucketName
	endpoint := cfg.Endpoint
	accessKeyID := cfg.AccessKeyID
	secretAccessKey := cfg.SecretAccessKey
	useSSL := cfg.UseSSL

	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}

	if err = client.MakeBucket(bucketName, "us-east-1"); err != nil {
		if err.Error() != bucketExistsErrMsg {
			return nil, err
		}
	}

	resource := fmt.Sprintf("arn:aws:s3:::%s/*", bucketName)

	policy := fmt.Sprintf(`{
		"Version":"2012-10-17",
		"Statement":[
		  {
			"Sid":"AddPerm",
			"Effect":"Allow",
			"Principal": "*",
			"Action": "s3:GetObject",
			"Resource": "%s"
		  }
		]
	  }`, resource)

	if err = client.SetBucketPolicy(bucketName, policy); err != nil {
		return nil, err
	}
	return client, nil
}
