package main

import (
	"fmt"
	"github.com/4726/discussion-board/services/posts/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error string
}

var (
	InvalidJSONBodyResponse = ErrorResponse{"invalid body"}
	PostDoesNotExist        = fmt.Errorf("post does not exist")
)

func CreatePost(db *gorm.DB, ctx *gin.Context) {
	form := struct {
		Title string `binding:"required"`
		Body  string `binding:"required"`
		User  string `binding:"required"`
	}{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.Set(logInfoKey, err)
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

	if err := db.Save(&post).Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"postID": post.ID})
}

func DeletePost(db *gorm.DB, ctx *gin.Context) {
	form := struct {
		PostID uint `binding:"required"`
		User   string
	}{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	post := models.Post{ID: form.PostID}

	if form.User != "" {
		if err := db.Where("user = ?", form.User).Delete(&post).Error; err != nil {
			ctx.Set(logInfoKey, err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		} else {
			ctx.JSON(http.StatusOK, gin.H{})
			return
		}
	}

	if err := db.Delete(&post).Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func UpdatePostLikes(db *gorm.DB, ctx *gin.Context) {
	form := struct {
		PostID uint `binding:"required"`
		Likes  int
	}{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	post := models.Post{ID: form.PostID}

	//uses UpdateColumn() instead of Update() because Update()
	//automatically updates the UpdatedAt field
	if err := db.Model(&post).UpdateColumn("Likes", form.Likes).Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func CreateComment(db *gorm.DB, ctx *gin.Context) {
	form := struct {
		PostID   uint `binding:"required"`
		ParentID uint
		User     string `binding:"required"`
		Body     string `binding:"required"`
	}{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.Set(logInfoKey, err)
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
			ctx.Set(logInfoKey, err)
			ctx.JSON(http.StatusBadRequest, ErrorResponse{"post does not exist"})
			return
		}
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func ClearComment(db *gorm.DB, ctx *gin.Context) {
	form := struct {
		CommentID uint `binding:"required"`
		User      string
	}{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	comment := models.Comment{ID: form.CommentID}

	if form.User != "" {
		if err := db.Model(&comment).Where("User = ?", form.User).UpdateColumn("Body", "").Error; err != nil {
			ctx.Set(logInfoKey, err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		} else {
			ctx.JSON(http.StatusOK, gin.H{})
		}
		return
	}

	if err := db.Model(&comment).UpdateColumn("Body", "").Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func UpdateCommentLikes(db *gorm.DB, ctx *gin.Context) {
	form := struct {
		CommentID uint `binding:"required"`
		Likes     int
	}{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	comment := models.Comment{ID: form.CommentID}

	if err := db.Model(&comment).Update("Likes", form.Likes).Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
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
