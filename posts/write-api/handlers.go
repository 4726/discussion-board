package main

import (
	"fmt"
	"github.com/4726/discussion-board/posts/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type CreateForm struct {
	Title, Body, User string
}

type DeleteForm struct {
	PostID uint
}

type UpdateLikesForm struct {
	PostID uint
	Likes  int
}

type CreateCommentForm struct {
	PostID, ParentID uint
	User, Body       string
}

type ClearCommentForm struct {
	CommentID uint
}

type UpdateCommentLikesForm struct {
	CommentID uint
	Likes     int
}

type ErrorResponse struct {
	Error string
}

var (
	InvalidJSONBodyResponse = ErrorResponse{"invalid body"}
	PostDoesNotExist        = fmt.Errorf("post does not exist")
)

func CreatePost(db *gorm.DB, ctx *gin.Context) {
	form := CreateForm{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	created := time.Now()
	post := models.Post{
		User:      form.User,
		Title:     form.Title,
		Body:      form.Body,
		Likes:     0,
		CreatedAt: created,
		UpdatedAt: created,
	}

	if err := validatePost(post); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	if err := db.Save(&post).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"postID": post.ID})
}

func DeletePost(db *gorm.DB, ctx *gin.Context) {
	form := DeleteForm{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	post := models.Post{ID: form.PostID}

	if err := db.Delete(&post).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func UpdatePostLikes(db *gorm.DB, ctx *gin.Context) {
	form := UpdateLikesForm{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	post := models.Post{ID: form.PostID}

	//uses UpdateColumn() instead of Update() because Update()
	//automatically updates the UpdatedAt field
	if err := db.Model(&post).UpdateColumn("Likes", form.Likes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func CreateComment(db *gorm.DB, ctx *gin.Context) {
	form := CreateCommentForm{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	created := time.Now()
	comment := models.Comment{
		PostID:    form.PostID,
		ParentID:  form.ParentID,
		User:      form.User,
		Body:      form.Body,
		CreatedAt: created,
		Likes:     0,
	}

	if err := addCommentToDB(db, &comment); err != nil {
		if err == PostDoesNotExist {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{"post does not exist"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func ClearComment(db *gorm.DB, ctx *gin.Context) {
	form := ClearCommentForm{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	comment := models.Comment{ID: form.CommentID}

	if err := db.Model(&comment).UpdateColumn("Body", "").Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func UpdateCommentLikes(db *gorm.DB, ctx *gin.Context) {
	form := UpdateCommentLikesForm{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	comment := models.Comment{ID: form.CommentID}

	if err := db.Model(&comment).Update("Likes", form.Likes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func deletePostFromDB(db *gorm.DB, postID uint) error {
	post := models.Post{ID: postID}
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
	if err := tx.Where("id = ?", postID).Delete("Comment{}").Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func addCommentToDB(db *gorm.DB, comment *models.Comment) error {
	post := models.Post{ID: comment.PostID}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	count := 0
	if err := tx.Model(&post).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}
	if count < 1 {
		tx.Rollback()
		return PostDoesNotExist
	}
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&post).Update("UpdatedAt", comment.CreatedAt).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func validatePost(p models.Post) error {
	if p.User == "" {
		return fmt.Errorf("empty user")
	}

	if p.Title == "" {
		return fmt.Errorf("empty title")
	}

	if p.Body == "" {
		return fmt.Errorf("empty body")
	}

	return nil
}
