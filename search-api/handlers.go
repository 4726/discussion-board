package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	Id    string
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
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	post := Post{form.Title, form.Body, form.User, form.Id, form.Timestamp, form.Likes}

	if err := esc.Index(post); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}

func Search(esc *ESClient, ctx *gin.Context) {
	term := ctx.Query("term")
	from := ctx.Query("from")
	total := ctx.Query("total")
	fromInt, err := strconv.Atoi(from)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid from query"})
		return
	}
	totalInt, err := strconv.Atoi(total)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid total query"})
		return
	}
	if err := verifySearchQuery(term, fromInt, totalInt); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	res, err := esc.Search(term, fromInt, totalInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func UpdateLikes(esc *ESClient, ctx *gin.Context) {
	var ulf UpdateLikesForm
	if err := ctx.BindJSON(&ulf); err != nil {
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}
	if err := esc.UpdateLikes(ulf.Id, ulf.Likes); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}

func Delete(esc *ESClient, ctx *gin.Context) {
	var df DeletePostForm
	if err := ctx.BindJSON(&df); err != nil {
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}
	if err := esc.Delete(df.Id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}

func UpdateLastUpdate(esc *ESClient, ctx *gin.Context) {
	var form UpdateLastUpdateForm
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}
	if err := esc.UpdateLastUpdate(form.Id, form.LastUpdate); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	ctx.JSON(http.StatusOK, struct{}{})
}

func verifySearchQuery(term string, from, total int) error {
	if total < 1 {
		return fmt.Errorf("invalid total query")
	}
	if term == "" {
		return fmt.Errorf("invalid term query")
	}
	return nil
}
