package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetPostMock(ctx *gin.Context) {
	postIDParam := ctx.Param("postid")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	data := map[string]interface{}{}
	data["id"] = postID
	data["user_id"] = 1
	data["title"] = "Hello World"
	data["body"] = "My first post"
	data["likes"] = 10
	data["created_at"] = time.Now().Truncate(time.Hour * 5)
	data["updated_at"] = time.Now().Truncate(time.Hour)
	data["has_like"] = true

	comment1 := map[string]interface{}{}
	comment1["id"] = 1
	comment1["post_id"] = postID
	comment1["parent_id"] = 0
	comment1["user_id"] = 2
	comment1["body"] = "good"
	comment1["created_at"] = data["created_at"].(time.Time).Add(time.Minute * 10)
	comment1["likes"] = 0
	comment1["has_like"] = false

	comment2 := map[string]interface{}{}
	comment2["id"] = 2
	comment2["post_id"] = postID
	comment2["parent_id"] = 0
	comment2["user_id"] = 3
	comment2["body"] = "great"
	comment2["created_at"] = data["created_at"].(time.Time).Add(time.Hour)
	comment2["likes"] = 1
	comment2["has_like"] = true

	comment3 := map[string]interface{}{}
	comment3["id"] = 3
	comment3["post_id"] = postID
	comment3["parent_id"] = 2
	comment3["user_id"] = 1
	comment3["body"] = "thank you"
	comment3["created_at"] = data["updated_at"]
	comment3["likes"] = 0
	comment3["has_like"] = false

	data["comments"] = []map[string]interface{}{comment1, comment2, comment3}
	ctx.JSON(http.StatusOK, data)
}

func GetPostsMock(ctx *gin.Context) {
	query := struct {
		Page   uint `form:"page" binding:"required"`
		UserID uint `form:"userid"`
	}{}
	if err := ctx.BindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	data1 := map[string]interface{}{}
	data1["id"] = 1
	data1["user_id"] = 1
	data1["title"] = "Hello World"
	data1["body"] = ""
	data1["likes"] = 10
	data1["created_at"] = time.Now().Truncate(time.Hour * 5)
	data1["updated_at"] = time.Now().Truncate(time.Hour)
	data1["comments"] = []gin.H{}

	posts := []gin.H{data1, data1, data1, data1, data1, data1, data1, data1, data1, data1}

	ctx.JSON(http.StatusOK, gin.H{"posts": posts})
}

func CreatePostMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"post_id": 1})
}

func DeletePostMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func LikePostMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"total": 11})
}

func UnlikePostMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"total": 11})
}

func AddCommentMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func LikeCommentMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"total": 11})
}

func UnlikeCommentMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"total": 11})
}

func ClearCommentMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func SearchMock(ctx *gin.Context) {
	form := struct {
		Term string `form:"term" binding:"required"`
		Page uint   `form:"page" binding:"required"`
	}{}
	if err := ctx.BindQuery(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	data1 := map[string]interface{}{}
	data1["id"] = 1
	data1["user_id"] = 1
	data1["title"] = "Hello World"
	data1["body"] = ""
	data1["Likes"] = 10
	data1["created_at"] = time.Now().Truncate(time.Hour * 5)
	data1["updated_at"] = time.Now().Truncate(time.Hour)
	data1["comments"] = []gin.H{}

	posts := []map[string]interface{}{data1, data1, data1, data1, data1, data1, data1, data1, data1, data1}

	ctx.JSON(http.StatusOK, gin.H{"posts": posts})
}

func RegisterGETMock(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"user_id": userID})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func RegisterPOSTMock(ctx *gin.Context) {
	if userID, err := getUserID(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"user_id": userID})
		return
	}
	jwt, err := generateJWT(1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"jwt": jwt, "user_id": 1})
}

func LoginGETMock(ctx *gin.Context) {
	if userID, err := getUserID(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"user_id": userID})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func LoginPOSTMock(ctx *gin.Context) {
	if userID, err := getUserID(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"user_id": userID})
		return
	}
	jwt, err := generateJWT(1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"jwt": jwt, "user_id": 1})
}

func ChangePasswordMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func GetProfileMock(ctx *gin.Context) {
	userIDParam := ctx.Param("userid")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	data := gin.H{}
	data["user_id"] = userID
	data["username"] = "my_username"
	data["bio"] = "This is my bio"
	data["avatar_id"] = ""
	if userID == 1 {
		data["is_mine"] = true
	}

	ctx.JSON(http.StatusOK, data)
}

func UpdateProfileMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func UserIdGETMock(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}

	ctx.JSON(http.StatusOK, gin.H{"user_id": userID})
}
