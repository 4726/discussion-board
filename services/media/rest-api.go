package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v6"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	bucketExistsErrMsg = "Your previous request to create the named bucket succeeded and you already own it."
	logInfoKey         = "log info"
)

var bucketName string

type RestAPI struct {
	mc     *minio.Client
	engine *gin.Engine
}

func NewRestAPI(cfg Config) (*RestAPI, error) {
	api := &RestAPI{}

	err := api.initMinio(cfg)
	if err != nil {
		return nil, err
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	api.engine = engine
	api.engine.Use(gin.Recovery())
	api.setRoutes()
	api.setMonitorRoute()
	api.engine.Use(log.RequestMiddleware())

	return api, nil
}

func (a *RestAPI) setRoutes() {
	a.engine.POST("/upload", func(ctx *gin.Context) {
		Upload(a.mc, ctx)
	})

	a.engine.POST("/remove/:name", func(ctx *gin.Context) {
		Remove(a.mc, ctx)
	})

	a.engine.GET("/info", func(ctx *gin.Context) {
		Info(a.mc, ctx)
	})
}

func (a *RestAPI) setMonitorRoute() {
	a.engine.Any("/metrics", gin.WrapH(promhttp.Handler()))
}

func (a *RestAPI) Run(addr string) error {
	return a.engine.Run(addr)
}

func (a *RestAPI) initMinio(cfg Config) error {
	bucketName = cfg.BucketName
	endpoint := cfg.Endpoint
	accessKeyID := cfg.AccessKeyID
	secretAccessKey := cfg.SecretAccessKey
	useSSL := cfg.UseSSL

	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return err
	}

	if err = client.MakeBucket(bucketName, "us-east-1"); err != nil {
		if err.Error() != bucketExistsErrMsg {
			return err
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
		return err
	}
	a.mc = client
	return nil
}
