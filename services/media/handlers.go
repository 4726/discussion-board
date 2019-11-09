package main

import (
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

type InfoReponse struct {
	StoreAddress string
}

func Upload(mc *minio.Client, ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("media")
	if err != nil {
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

	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	guid, err := ksuid.NewRandom() //not guaranteed unique
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	name := guid.String()
	_, err = mc.PutObject(bucketName, name, file, fileHeader.Size, minio.PutObjectOptions{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, UploadResponse{name})
}

func Remove(mc *minio.Client, ctx *gin.Context) {
	name := ctx.Param("name")
	if err := mc.RemoveObject(bucketName, name); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, struct{}{})
}

func Info(mc *minio.Client, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, InfoReponse{mc.EndpointURL().String()})
}
