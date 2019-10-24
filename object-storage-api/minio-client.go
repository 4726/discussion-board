package main

import (
	"github.com/minio/minio-go/v6"
	"github.com/segmentio/ksuid"
	"io"
)

type MinioClient struct {
	client   *minio.Client
	Endpoint string
}

const (
	BUCKET_EXISTS = "Your previous request to create the named bucket succeeded and you already own it."
)

func NewMinioClient(cfg Config) (*MinioClient, error) {
	endpoint := cfg.Endpoint
	accessKeyID := cfg.AccessKeyID
	secretAccessKey := cfg.SecretAccessKey
	useSSL := cfg.UseSSL

	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}

	if err = client.MakeBucket("images", "us-east-1"); err != nil {
		if err.Error() != BUCKET_EXISTS {
			return nil, err
		}
	}

	policy := `{
		"Version":"2012-10-17",
		"Statement":[
		  {
			"Sid":"AddPerm",
			"Effect":"Allow",
			"Principal": "*",
			"Action": "s3:GetObject",
			"Resource": "arn:aws:s3:::images/*"
		  }
		]
	  }`

	if err = client.SetBucketPolicy("images", policy); err != nil {
		return nil, err
	}

	return &MinioClient{client, endpoint}, nil
}

func (mc *MinioClient) PutImage(r io.Reader, size int64) (string, error) {
	guid, err := ksuid.NewRandom() //not guaranteed unique
	if err != nil {
		return "", err
	}
	name := guid.String()
	_, err = mc.client.PutObject("images", name, r, size, minio.PutObjectOptions{})
	return name, err
}

func (mc *MinioClient) PutImageFile(file string) (string, error) {
	guid, err := ksuid.NewRandom() //not guaranteed unique
	if err != nil {
		return "", err
	}
	name := guid.String()
	_, err = mc.client.FPutObject("images", name, file, minio.PutObjectOptions{})
	return name, err
}

func (mc *MinioClient) RemoveImage(name string) error {
	return mc.client.RemoveObject("images", name)
}
