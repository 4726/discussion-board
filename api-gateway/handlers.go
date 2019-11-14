package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var jwtSecretKey = []byte("todosecretkey")

type JWTClaims struct {
	jwt.StandardClaims
	UserID uint
}

func GetPost(ctx *gin.Context) {
	postIDParam := ctx.Param("postid")
	resp, _ := get(PostsReadServiceAddr() + "/posts/" + postIDParam)
	ctx.JSON(resp.StatusCode, resp.Data)
}

func GetPosts(ctx *gin.Context) {
	resp, _ := get(fmt.Sprintf("%s/posts?from=%v&total=%v&user=%v&sort=%v",
		PostsReadServiceAddr(), 0, 10, "", ""))
	ctx.JSON(resp.StatusCode, resp.Data)
}

func CreatePost(ctx *gin.Context) {
	defer ctx.Request.Body.Close()
	resp, _ := postProxy(PostsWriteServiceAddr()+"/post/create", ctx.Request.Body)
	ctx.JSON(resp.StatusCode, resp.Data)
}

func DeletePost(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	resp, _ := post(PostsWriteServiceAddr()+"/post/delete", m)
	ctx.JSON(resp.StatusCode, resp.Data)
}

func LikePost(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	resp, _ := post(LikesServiceAddr()+"/post/like", m)

	//should not matter much if it fails since
	//it can try again when another user likes/unlikes
	go func() {
		data := struct {
			PostID uint
			Likes  int
		}{m["PostID"].(uint), resp.Data["Total"].(int)}
		_, _ = post(PostsWriteServiceAddr()+"/post/likes", data)
	}()

	ctx.JSON(resp.StatusCode, resp.Data)
}

func UnlikePost(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	resp, _ := post(LikesServiceAddr()+"/post/unlike", m)

	go func() {
		data := struct {
			PostID uint
			Likes  int
		}{m["PostID"].(uint), resp.Data["Total"].(int)}
		_, _ = post(PostsWriteServiceAddr()+"/post/likes", data)
	}()

	ctx.JSON(resp.StatusCode, resp.Data)
}

func AddComment(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	resp, _ := post(PostsWriteServiceAddr()+"/comment/create", m)
	ctx.JSON(resp.StatusCode, resp.Data)
}

func LikeComment(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	resp, _ := post(LikesServiceAddr()+"/comment/like", m)

	go func() {
		data := struct {
			PostID uint
			Likes  int
		}{m["PostID"].(uint), resp.Data["Total"].(int)}
		_, _ = post(PostsWriteServiceAddr()+"/post/likes", data)
	}()

	ctx.JSON(resp.StatusCode, resp.Data)
}

func UnlikeComment(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	resp, _ := post(LikesServiceAddr()+"/post/unlike", m)

	go func() {
		data := struct {
			PostID uint
			Likes  int
		}{m["PostID"].(uint), resp.Data["Total"].(int)}
		_, _ = post(PostsWriteServiceAddr()+"/post/likes", data)
	}()

	ctx.JSON(resp.StatusCode, resp.Data)
}

func ClearComment(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	postIDParam := ctx.Param("postid")
	resp, _ := post(PostsWriteServiceAddr()+"/comment/clear/"+postIDParam, m)
	ctx.JSON(resp.StatusCode, resp.Data)
}

func Search(ctx *gin.Context) {
	form := struct {
		Term string `form:"term" binding:"required"`
		Page uint   `form:"page" binding:"required"`
	}{}
	if err := ctx.BindQuery(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	resp, _ := get(fmt.Sprintf("%s/search?from=%v&total=%v&term=%v",
		SearchServiceAddr(), (form.Page*10)-10, 10, form.Term))
	ctx.JSON(resp.StatusCode, resp.Data)
}

func RegisterGET(ctx *gin.Context) {
	if _, err := getUserID(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func RegisterPOST(ctx *gin.Context) {
	if _, err := getUserID(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	defer ctx.Request.Body.Close()
	resp, _ := postProxy(UserServiceAddr()+"/account", ctx.Request.Body)
	jwt, err := generateJWT(resp.Data["userID"].(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	resp.Data["jwt"] = jwt
	ctx.JSON(resp.StatusCode, resp.Data)
}

func LoginGET(ctx *gin.Context) {
	if _, err := getUserID(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func LoginPOST(ctx *gin.Context) {
	if _, err := getUserID(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	defer ctx.Request.Body.Close()
	resp, _ := postProxy(UserServiceAddr()+"/login", ctx.Request.Body)
	jwt, err := generateJWT(resp.Data["userID"].(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	resp.Data["jwt"] = jwt
	ctx.JSON(resp.StatusCode, resp.Data)
}

func ChangePassword(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	resp, _ := post(UserServiceAddr()+"/password", m)
	ctx.JSON(resp.StatusCode, resp.Data)
}

func GetProfile(ctx *gin.Context) {
	userIDParam := ctx.Param("userid")
	resp, _ := get(UserServiceAddr() + "/" + userIDParam)
	ctx.JSON(resp.StatusCode, resp.Data)
}

func UpdateProfile(ctx *gin.Context) {
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

	m, err := bindJSONAndAdd(ctx, extra)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	_, _ = post(UserServiceAddr()+"/profile/update", m)

	ctx.JSON(http.StatusOK, gin.H{})
}

func generateJWT(userID uint) (string, error) {
	claims := JWTClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

//jwt must be stored in authorization bearer
//format: Authorization: Bearer <token>
func getUserID(ctx *gin.Context) (uint, error) {
	authHeader := ctx.GetHeader("Authorization")
	splitTokens := strings.Split(authHeader, "Bearer ")
	if len(splitTokens) != 2 {
		return 0, fmt.Errorf("invalid authorization header format")
	}
	reqToken := strings.TrimSpace(splitTokens[1])
	token, err := jwt.ParseWithClaims(reqToken, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signingmethod")
		}

		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		userID := claims.UserID
		if userID < 1 {
			return 0, fmt.Errorf("invalid userid")
		}
		return userID, nil
	}

	//should not happen
	return 0, fmt.Errorf("wrong claims type")
}
