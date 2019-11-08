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
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		ctx.Set(logInfoKey, err)
		if gorm.IsRecordNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, post)
}

func GetPosts(db *gorm.DB, ctx *gin.Context) {
	totalS := ctx.Query("total")
	fromS := ctx.Query("from")
	user := ctx.Query("user")
	sortType := ctx.Query("sort")

	total, err := strconv.Atoi(totalS)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid total query"})
		return
	}
	from, err := strconv.Atoi(fromS)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid from query"})
		return
	}

	switch sortType {
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

	if user != "" {
		posts, err := getPostsUser(db, from, total, user, sortType)
		if err != nil {
			ctx.Set(logInfoKey, err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, posts)
		return
	}

	posts, err := getPosts(db, from, total, sortType)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func getPosts(db *gorm.DB, from, total int, sortType string) ([]models.Post, error) {
	posts := []models.Post{}
	selectFields := []string{"id", "user", "title", "likes", "created_at", "updated_at"}
	if err := db.Preload("Comments").Select(selectFields).
		Order(sortType).
		Offset(from).
		Limit(total).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func getPostsUser(db *gorm.DB, from, total int, user, sortType string) ([]models.Post, error) {
	posts := []models.Post{}
	selectFields := []string{"id", "user", "title", "likes", "created_at", "updated_at"}
	if err := db.Preload("Comments").Select(selectFields).
		Where("user = ?", user).
		Order(sortType).
		Offset(from).
		Limit(total).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}
