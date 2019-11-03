package main

import (
	"github.com/4726/discussion-board/posts/models"
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
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
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
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid total query"})
		return
	}
	from, err := strconv.Atoi(fromS)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid from query"})
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
		return
	}

	posts, err := getPosts(db, from, total, sortType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

//what happens to non selected fields? set to zero value?
func getPosts(db *gorm.DB, from, total int, sortType string) ([]models.Post, error) {
	posts := []models.Post{}
	selectFields := []string{"Post_ID", "User", "Title", "Likes", "Created_At", "Updated_At"}
	if err := db.Select(selectFields).
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
	selectFields := []string{"Post_ID", "User", "Title", "Likes", "Created_At", "Updated_At"}
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
