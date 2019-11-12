package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type PostLikeForm struct {
	PostID, UserID uint
}

type CommentLikeForm struct {
	CommentID, UserID uint
}

type IDsForm struct {
	IDs []uint
}

type ErrorResponse struct {
	Error string
}

type IDLikes struct {
	ID    uint
	Likes int
}

//should be fine without transaction
func LikePost(db *gorm.DB, ctx *gin.Context) {
	form := &PostLikeForm{}
	if err := ctx.BindJSON(form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid form"})
		return
	}

	like := PostLike{form.PostID, form.UserID, time.Now()}

	if err := db.FirstOrCreate(&PostLike{}, &like).Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	var count int

	if err := db.Where("post_id = ?", like.PostID).Find(&PostLike{}).Count(&count).Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"total": count})
}

func UnlikePost(db *gorm.DB, ctx *gin.Context) {
	form := &PostLikeForm{}
	if err := ctx.BindJSON(form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid form"})
		return
	}

	like := PostLike{PostID: form.PostID, UserID: form.UserID}

	if err := db.Delete(&like).Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	var count int

	if err := db.Where("post_id = ?", like.PostID).Find(&PostLike{}).Count(&count).Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"total": count})
}

func LikeComment(db *gorm.DB, ctx *gin.Context) {
	form := &CommentLikeForm{}
	if err := ctx.BindJSON(form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid form"})
		return
	}

	like := CommentLike{form.CommentID, form.UserID, time.Now()}

	if err := db.FirstOrCreate(&CommentLike{}, &like).Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	var count int

	if err := db.Where("comment_id = ?", like.CommentID).Find(&CommentLike{}).Count(&count).Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"total": count})
}

func UnlikeComment(db *gorm.DB, ctx *gin.Context) {
	form := &CommentLikeForm{}
	if err := ctx.BindJSON(form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid form"})
		return
	}

	like := CommentLike{CommentID: form.CommentID, UserID: form.UserID}

	if err := db.Delete(&like).Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	var count int

	if err := db.Where("comment_id = ?", like.CommentID).Find(&CommentLike{}).Count(&count).Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"total": count})
}

func GetMultiplePostLikes(db *gorm.DB, ctx *gin.Context) {
	form := &IDsForm{}
	if err := ctx.BindJSON(form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid form"})
		return
	}

	counts := []IDLikes{}

	for _, v := range form.IDs {
		var count int

		if err := db.Where("post_id = ?", v).Find(&PostLike{}).Count(&count).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				counts = append(counts, IDLikes{v, 0})
				continue
			}
			ctx.Set(logInfoKey, err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
			return
		}

		counts = append(counts, IDLikes{v, count})
	}

	ctx.JSON(http.StatusOK, counts)
}

func GetMultipleCommentLikes(db *gorm.DB, ctx *gin.Context) {
	form := &IDsForm{}
	if err := ctx.BindJSON(form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid form"})
		return
	}

	counts := []IDLikes{}

	for _, v := range form.IDs {
		var count int

		if err := db.Where("comment_id = ?", v).Find(&CommentLike{}).Count(&count).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				counts = append(counts, IDLikes{v, 0})
				continue
			}
			ctx.Set(logInfoKey, err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
			return
		}

		counts = append(counts, IDLikes{v, count})
	}

	ctx.JSON(http.StatusOK, counts)
}
