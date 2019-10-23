package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type IndexForm struct {
	Title, Body, User, Id string
}

type UpdateLikesForm struct {
	Id    string
	Likes int
}

type ErrorResponse struct {
	Error string
}

func Search(esc *ESClient, ctx *gin.Context) {
	term := ctx.Query("term")
	from := ctx.Query("from")
	total := ctx.Query("total")
	fromInt, err := strconv.Atoi(from)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"from is not an int"})
		return
	}
	totalInt, err := strconv.Atoi(total)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"total is not an int"})
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

func Index(esc *ESClient, ctx *gin.Context) {
	var iform IndexForm
	if err := ctx.BindJSON(&iform); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	if err := esc.Index(iform.Title, iform.Body, iform.User, iform.Id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, struct{}{})
}

func UpdateLikes(esc *ESClient, ctx *gin.Context) {
	var ulf UpdateLikesForm
	if err := ctx.BindJSON(&ulf); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	if err := esc.UpdateLikes(ulf.Id, ulf.Likes); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, struct{}{})
}

func Delete(esc *ESClient, ctx *gin.Context) {
	id := ctx.Param("id")
	if err := esc.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, struct{}{})
}

func verifySearchQuery(term string, from, total int) error {
	if term == "" {
		return fmt.Errorf("term cannot be empty")
	}
	if total == 0 {
		return fmt.Errorf("total cannot be 0")
	}
	return nil
}
