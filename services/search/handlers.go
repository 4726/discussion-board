package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexForm struct {
	Title, Body, User, Id string
	Timestamp             int64
	Likes                 int
}

type UpdateLikesForm struct {
	Id    string
	Likes int
}

type DeletePostForm struct {
	Id string
}

type UpdateLastUpdateForm struct {
	Id         string
	LastUpdate int64
}

type ErrorResponse struct {
	Error string
}

var (
	InvalidJSONBodyResponse = ErrorResponse{"invalid body"}
)

func Index(esc *ESClient, ctx *gin.Context) {
	var form IndexForm
	if err := ctx.BindJSON(&form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	post := Post{form.Title, form.Body, form.User, form.Id, form.Timestamp, form.Likes}

	if err := esc.Index(post); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}

func Search(esc *ESClient, ctx *gin.Context) {
	query := struct{
		Term string `form:"term" binding:"required"`
		From uint `form:"from"`
		Total uint `form:"total" binding:"required"`
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
	var ulf UpdateLikesForm
	if err := ctx.BindJSON(&ulf); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}
	if err := esc.UpdateLikes(ulf.Id, ulf.Likes); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}

func Delete(esc *ESClient, ctx *gin.Context) {
	var df DeletePostForm
	if err := ctx.BindJSON(&df); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}
	if err := esc.Delete(df.Id); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}

func UpdateLastUpdate(esc *ESClient, ctx *gin.Context) {
	var form UpdateLastUpdateForm
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