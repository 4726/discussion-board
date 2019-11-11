package main

import (
	"strconv"
	"net/http"
)

func GetPost(ctx *gin.Context) {
	postIDParam := ctx.Param("postid")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	post, err := postsService.Get(postID)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	ctx.HTML(http.StatusOK, "post.html", gin.H{"post": post})
}

func GetPosts(ctx *gin.Context) {
	posts, err := postsService.GetMany(10, 0)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	ctx.HTML(http.StatusOK, "post.html", gin.H{"posts": posts})
}

func DeletePost(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.Redirect(http.StatusUnauthorized, "home.html", nil)
		return
	}
	postIDParam := ctx.Param("postid")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		ctx.HTML(http.StatusNoContent, "404.html", nil)
		return
	}
	if err := postsService.DeleteIfOwner(postID, userID); err != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	ctx.HTML(http.StatusNoContent, "post.html", nil)
}

func LikePost(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.Redirect(http.StatusUnauthorized, "home.html", nil)
		return
	}
	postIDParam := ctx.Param("postid")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	if err := postsService.LikePost(postID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func UnlikePost(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.Redirect(http.StatusUnauthorized, "home.html", nil)
		return
	}
	postIDParam := ctx.Param("postid")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	if err := postsService.UnlikePost(postID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}


func AddComment(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.Redirect(http.StatusUnauthorized, "home.html", nil)
		return
	}
	postIDForm := ctx.PostForm("postID")
	parentIDForm := ctx.PostForm("parentID")
	body := ctx.PostForm("body")
	if postID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	postID, err := strconv.Atoi(postIDForm)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	parentID, err := strconv.Atoi(parentIDForm)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	if err := postsService.CreateComment(CommentData{postID, parentID, userID, body}); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func LikeComment(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.Redirect(http.StatusUnauthorized, "home.html", nil)
		return
	}
	commentIDParam := ctx.Param("commentid")
	commentID, err := strconv.Atoi(postIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	if err := postsService.LikeComment(commentID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func UnlikeComment(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.Redirect(http.StatusUnauthorized, "home.html", nil)
		return
	}
	commentIDParam := ctx.Param("commentid")
	commentID, err := strconv.Atoi(postIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	if err := postsService.UnlikeComment(commentID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func ClearComment(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.Redirect(http.StatusUnauthorized, "home.html", nil)
		return
	}
	commentIDParam := ctx.Param("commentid")
	commentID, err := strconv.Atoi(commentIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	if err := postsService.DeleteIfOwner(commentID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func Search(ctx *gin.Context) {
	term := ctx.Query("term")
	if term == "" {
		ctx.HTML(http.BadRequest, "400.html", nil)
		return
	}
	pageParam := ctx.Query("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		ctx.HTML(http.BadRequest, "400.html", nil)
		return
	}
	postIDs, err := searchService.Search(10 * page, 10, term)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	posts, err := postsService.GetMultiple(postIDs)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	ctx.HTML(http.StatusOK, "posts.html", gin.H{"posts": posts})
}

func RegisterGET(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID != 0 {
		ctx.Redirect(http.StatusOK, "home.html", nil)
		return
	}

	ctx.HTML(http.StatusOK, "register.html", nil)
}

func RegisterPOST(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID != 0 {
		ctx.Redirect(http.StatusOK, "home.html", nil)
		return
	}

	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	userID, err := userService.CreateAccount(username, password)
	if err != nil {
		switch error.Error() {
		case "invalid username":
			ctx.HTML(http.StatusBadRequest, "400.html", nil)
		case "invalid password":
			ctx.HTML(http.StatusBadRequest, "400.html", nil)
		case "username unavailable":
			ctx.HTML(http.StatusBadRequest, "400.html", nil)
		default:
			ctx.HTML(http.StatusInternalServerError, "500.html", nil)
		}
		return
	}

	//store jwt

	ctx.HTML(http.StatusOK, "home.html", nil)
}

func LoginGET(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID != 0 {
		ctx.Redirect(http.StatusOK, "home.html", nil)
		return
	}

	ctx.HTML(http.StatusOK, "login.html", nil)
}

func LoginPOST(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID != 0 {
		ctx.Redirect(http.StatusOK, "home.html", nil)
		return
	}

	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	userID, err := userService.ValidLogin(username, password)
	if err != nil {
		switch error.Error() {
		case "invalid login":
			ctx.HTML(http.StatusUnauthorized, "401.html", nil)
		default:
			ctx.HTML(http.StatusInternalServerError, "500.html", nil)
		}
		return
	}

	//store jwt 

	ctx.HTML(http.StatusOK, "home.html", nil)
}

func ChangePassword(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.Redirect(http.StatusUnauthorized, "home.html", nil)
		return
	}

	username := ctx.PostForm("username")
	oldPass := ctx.PostForm("oldpass")
	newPass := ctx.PostForm("newpass")

	err := userService.ChangePassword(username, oldPass, newPass)
	if err != nil {
		switch error.Error() {
		case "invalid old password":
			ctx.HTML(http.StatusBadRequest, "400.html", nil)
		case "invalid new password":
			ctx.HTML(http.StatusBadRequest, "400.html", nil)
		default:
			ctx.HTML(http.StatusInternalServerError, "500.html", nil)
		}
		return
	}

	ctx.HTML(http.StatusOK, "home.html", nil)
}

func GetProfile(ctx *gin.Context) {
	userIDParam := ctx.Param("userid")
	userID, err := strconv.Atoi(postIDParam)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	profile, err := userService.GetProfile(userID)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	ctx.HTML(http.StatusOK, "profile.html", gin.H{"profile": profile})
}

func UpdateProfile(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.Redirect(http.StatusUnauthorized, "home.html", nil)
		return
	}

	opts := UpdateProfileOptions{}

	bio := ctx.PostForm("bio")
	if bio != "" {
		opts.Bio = bio
	}
	fileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "400.html", nil)
		return
	} else {
		objectName, err := mediaService.Upload(fileHeader.Filename)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}
		opts.Avatar = objectName
	}

	err := userService.UpdateProfile(userID, opts)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	ctx.HTML(http.StatusOK, "home.html", nil)
}

//probably use jwt
func getUserID(ctx *gin.Context) int {
	return 0
}