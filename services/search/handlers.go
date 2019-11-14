package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResponse struct {
	Error string
}

var (
	InvalidJSONBodyResponse = ErrorResponse{"invalid body"}
)

func Index(esc *ESClient, ctx *gin.Context) {
	form := struct {
		Title     string `binding:"required"`
		Body      string `binding:"required"`
		UserID    uint   `binding:"required"`
		Id        string `binding:"required"`
		Timestamp int64
		Likes     int
	}{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	post := Post{form.Title, form.Body, form.Id, int(form.UserID), form.Timestamp, form.Likes}

	if err := esc.Index(post); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}

func Search(esc *ESClient, ctx *gin.Context) {
	query := struct {
		Term  string `form:"term" binding:"required"`
		From  uint   `form:"from"`
		Total uint   `form:"total" binding:"required"`
	}{}

	err := ctx.BindQuery(&query)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid query"})
		return
	}

	res, err := esc.Search(query.Term, int(query.From), int(query.Total))
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func UpdateLikes(esc *ESClient, ctx *gin.Context) {
	form := struct {
		Id    string `binding:"required"`
		Likes int
	}{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}
	if err := esc.UpdateLikes(form.Id, form.Likes); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}

func Delete(esc *ESClient, ctx *gin.Context) {
	form := struct {
		Id string `binding:"required"`
	}{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}
	if err := esc.Delete(form.Id); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}

func UpdateLastUpdate(esc *ESClient, ctx *gin.Context) {
	form := struct {
		Id         string `binding:"required"`
		LastUpdate int64
	}{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}
	if err := esc.UpdateLastUpdate(form.Id, form.LastUpdate); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}
