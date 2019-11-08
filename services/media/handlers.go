package main

import (
	"github.com/gin-gonic/gin"
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

func Upload(mc *MinioClient, ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	name, err := mc.PutImage(file, fileHeader.Size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, UploadResponse{name})
}

func Remove(mc *MinioClient, ctx *gin.Context) {
	name := ctx.Param("name")
	if err := mc.RemoveImage(name); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, struct{}{})
}

func Info(mc *MinioClient, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, InfoReponse{mc.Endpoint})
}
