package main

import (
	"net/http"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPostMock(ctx *gin.Context) {
	postIDParam := ctx.Param("postid")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}
	data := map[string]interface{}{}
	data["ID"] = postID
	data["UserID"] = 1
	data["Title"] = "Hello World"
	data["Body"] = "My first post"
	data["Likes"] = 10
	data["CreatedAt"] = time.Now().Truncate(time.Hour * 5)
	data["UpdatedAt"] = time.Now().Truncate(time.Hour)

	comment1 := map[string]interface{}{}
	comment1["ID"] = 1
	comment1["PostID"] = postID
	comment1["ParentID"] = 0
	comment1["UserID"] = 2
	comment1["Body"] = "good"
	comment1["CreatedAt"] = data["CreatedAt"].(time.Time).Add(time.Minute * 10)
	comment1["Likes"] = 0

	comment2 := map[string]interface{}{}
	comment2["ID"] = 2
	comment2["PostID"] = postID
	comment2["ParentID"] = 0
	comment2["UserID"] = 3
	comment2["Body"] = "great"
	comment2["CreatedAt"] = data["CreatedAt"].(time.Time).Add(time.Hour)
	comment2["Likes"] = 1

	comment3 := map[string]interface{}{}
	comment3["ID"] = 3
	comment3["PostID"] = postID
	comment3["ParentID"] = 2
	comment3["UserID"] = 1
	comment3["Body"] = "thank you"
	comment3["CreatedAt"] = data["UpdatedAt"]
	comment3["Likes"] = 0

	data["Comments"] = []map[string]interface{}{comment1, comment2, comment3}
	ctx.JSON(http.StatusOK, data)
}

func GetPostsMock(ctx *gin.Context) {
	query := struct {
		Page  uint   `form:"page" binding:"required"`
		UserID uint   `form:"userid"`
	}{}
	if err := ctx.BindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	data1 := map[string]interface{}{}
	data1["ID"] = 1
	data1["UserID"] = 1
	data1["Title"] = "Hello World"
	data1["Body"] = ""
	data1["Likes"] = 10
	data1["CreatedAt"] = time.Now().Truncate(time.Hour * 5)
	data1["UpdatedAt"] = time.Now().Truncate(time.Hour)
	data1["Comments"] = []gin.H{}
	
	ctx.JSON(http.StatusOK, []gin.H{data1, data1, data1, data1, data1, data1, data1, data1, data1, data1})
}

func CreatePostMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"postID": 1})
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

	ctx.JSON(http.StatusOK, gin.H{})
}

func UnlikePostMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
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

	ctx.JSON(http.StatusOK, gin.H{})
}

func UnlikeCommentMock(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
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
	data1["ID"] = 1
	data1["UserID"] = 1
	data1["Title"] = "Hello World"
	data1["Body"] = ""
	data1["Likes"] = 10
	data1["CreatedAt"] = time.Now().Truncate(time.Hour * 5)
	data1["UpdatedAt"] = time.Now().Truncate(time.Hour)
	data1["Comments"] = []gin.H{}
	
	ctx.JSON(http.StatusOK, []gin.H{data1, data1, data1, data1, data1, data1, data1, data1, data1, data1})
}

func RegisterGETMock(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"UserID": userID})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func RegisterPOSTMock(ctx *gin.Context) {
	if userID, err := getUserID(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"userID": userID})
		return
	}
	jwt, err := generateJWT(1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"jwt": jwt, "userID": 1})
}

func LoginGETMock(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"UserID": userID})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func LoginPOSTMock(ctx *gin.Context) {
	if userID, err := getUserID(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"userID": userID})
		return
	}
	jwt, err := generateJWT(1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"jwt": jwt, "userID": 1})
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
	data["UserID"] = userID
	data["Username"] = "my_username"
	data["Bio"] = "This is my bio"
	data["AvatarID"] = ""

	ctx.JSON(http.StatusOK, data)
}

func UpdateProfileMock(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	extra := gin.H{"UserID": userID}

	_, err = ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	} else {
		//upload multipart form to media service
		//then get the avatar id
		//then add to extra map
	}

	newBody, ok := ctx.GetPostForm("body")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	extra["body"] = newBody

	ctx.JSON(http.StatusOK, gin.H{})
}