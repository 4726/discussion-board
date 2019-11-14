package main

import (
	"github.com/4726/discussion-board/services/posts/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Error string
}

func GetFullPost(db *gorm.DB, ctx *gin.Context) {
	postIDS := ctx.Param("postid")

	postID, err := strconv.Atoi(postIDS)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusNotFound, struct{}{})
		return
	}

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		ctx.Set(logInfoKey, err)
		if gorm.IsRecordNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, struct{}{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, post)
}

func GetPosts(db *gorm.DB, ctx *gin.Context) {
	query := struct {
		Total  uint   `form:"total" binding:"required"`
		From   uint   `form:"from"`
		UserID uint   `form:"userid"`
		Sort   string `form:"sort"`
	}{}

	if err := ctx.BindQuery(&query); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid query"})
		return
	}

	var sortType string

	switch query.Sort {
	case "likes_desc":
		sortType = "likes desc"
	case "created_at_desc":
		sortType = "created_at desc"
	case "created_at":
		sortType = "created_at"
	case "updated_at":
		sortType = "updated_at"
	default:
		sortType = "updated_at desc"
	}

	if query.UserID != 0 {
		posts, err := getPostsUser(db, query.From, query.Total, query.UserID, sortType)
		if err != nil {
			ctx.Set(logInfoKey, err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
			return
		}

		ctx.JSON(http.StatusOK, posts)
		return
	}

	posts, err := getPosts(db, query.From, query.Total, sortType)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func getPosts(db *gorm.DB, from, total uint, sortType string) ([]models.Post, error) {
	posts := []models.Post{}
	selectFields := []string{"id", "user_id", "title", "likes", "created_at", "updated_at"}
	if err := db.Preload("Comments").Select(selectFields).
		Order(sortType).
		Offset(from).
		Limit(total).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func getPostsUser(db *gorm.DB, from, total uint, userID uint, sortType string) ([]models.Post, error) {
	posts := []models.Post{}
	selectFields := []string{"id", "user_id", "title", "likes", "created_at", "updated_at"}
	if err := db.Preload("Comments").Select(selectFields).
		Where("user_id = ?", userID).
		Order(sortType).
		Offset(from).
		Limit(total).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}
