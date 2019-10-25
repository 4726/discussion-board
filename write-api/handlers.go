package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type CreateForm struct {
	Title, Body, User string
}

type DeleteForm struct {
	PostID int
}

type UpdateLikesForm struct {
	PostID, Likes int
}

type AddCommentForm struct {
	PostID, ParentID int
	User, Body       string
}

type ClearCommentForm struct {
	CommentID int
}

type UpdateCommentLikesForm struct {
	CommentID, Likes int
}

type ErrorResponse struct {
	Error string
}

func CreatePost(db *gorm.DB, ctx *gin.Context) {
	form := CreateForm{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	created := time.Now()
	post := Post{
		User:      form.User,
		Title:     form.Title,
		Body:      form.Body,
		Likes:     0,
		CreatedAt: created,
		UpdatedAt: created,
	}

	if err := db.Create(&post).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func DeletePost(db *gorm.DB, ctx *gin.Context) {
	form := DeleteForm{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	if err := deletePostFromDB(db, form.PostID); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func UpdatePostLikes(db *gorm.DB, ctx *gin.Context) {
	form := UpdateLikesForm{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	post := Post{PostID: form.PostID}

	if err := db.Model(&post).Update("Likes", form.Likes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func AddComment(db *gorm.DB, ctx *gin.Context) {
	form := AddCommentForm{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	created := time.Now()
	comment := Comment{
		PostID:    form.PostID,
		ParentID:  form.ParentID,
		User:      form.User,
		Body:      form.Body,
		CreatedAt: created,
		Likes:     0,
	}

	if err := addCommentToDB(db, &comment); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func ClearComment(db *gorm.DB, ctx *gin.Context) {
	form := ClearCommentForm{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	comment := Comment{CommentID: form.CommentID}

	if err := db.Model(&comment).Update("Body", "").Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func UpdateCommentLikes(db *gorm.DB, ctx *gin.Context) {
	form := UpdateCommentLikesForm{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	comment := Comment{CommentID: form.CommentID}

	if err := db.Model(&comment).Update("Likes", form.Likes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func deletePostFromDB(db *gorm.DB, postID int) error {
	post := Post{PostID: postID}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("PostID = ?", postID).Delete("Comment{}").Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func addCommentToDB(db *gorm.DB, comment *Comment) error {
	post := Post{PostID: comment.PostID}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&post).Update("Updated", comment.CreatedAt).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
