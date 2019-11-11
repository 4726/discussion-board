package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v6"
	"github.com/segmentio/ksuid"
	"net/http"
)

type ErrorResponse struct {
	Error string
}

type UploadResponse struct {
	Name string
}

type InfoResponse struct {
	StoreAddress string
}

func Upload(mc *minio.Client, ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("media")
	if err != nil {
		ctx.Set(logInfoKey, err)
		if err.Error() == "missing form body" {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid form"})
			return
		}
		if err.Error() == "http: no such file" {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid form"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	contentTypes, ok := fileHeader.Header["Content-Type"]
	if !ok {
		ctx.Set(logInfoKey, "no Content-Type in header")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	if len(contentTypes) == 0 {
		ctx.Set(logInfoKey, "no Content-Type value set")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	contentType := contentTypes[0]

	file, err := fileHeader.Open()
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	defer file.Close()

	guid, err := ksuid.NewRandom() //not guaranteed unique
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	name := guid.String()
	opts := minio.PutObjectOptions{ContentType: contentType}
	_, err = mc.PutObject(bucketName, name, file, fileHeader.Size, opts)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, UploadResponse{name})
}

func Remove(mc *minio.Client, ctx *gin.Context) {
	name := ctx.Param("name")
	if err := mc.RemoveObject(bucketName, name); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}

func Info(mc *minio.Client, ctx *gin.Context) {
	endpoint := fmt.Sprintf("%s/%s/", mc.EndpointURL().String(), bucketName)
	ctx.JSON(http.StatusOK, InfoResponse{endpoint})
}
