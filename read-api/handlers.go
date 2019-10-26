package main

import (
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
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
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
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	from, err := strconv.Atoi(fromS)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	switch sortType {
	case "likes":
		sortType = "Likes"
	case "createdat":
		sortType = "Created_At"
	default:
		sortType = "Updated_At"
	}

	if user != "" {
		posts, err := getPostsUser(db, from, total, user, sortType)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, posts)
	}

	posts, err := getPosts(db, from, total, sortType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

//what happens to non selected fields? set to zero value?
func getPosts(db *gorm.DB, from, total int, sortType string) ([]Post, error) {
	posts := []Post{}
	selectFields := []string{"Post_ID", "Likes", "User", "Title", "Created_At", "Updated_At"}
	if err := db.Select(selectFields).
		Order(sortType).
		Offset(from).
		Limit(total).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func getPostsUser(db *gorm.DB, from, total int, user, sortType string) ([]Post, error) {
	posts := []Post{}
	selectFields := []string{"Post_ID", "Likes", "User", "Title", "Created_At", "Updated_At"}
	if err := db.Select(selectFields).
		Where("user = ?", user).
		Order(sortType).
		Offset(from).
		Limit(total).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}
