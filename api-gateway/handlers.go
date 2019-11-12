package main

import (
	"strconv"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"fmt"
)

var jwtSecretKey = []byte("todosecretkey")

type ErrorResponse{
	Error string
}

type JWTClaims struct {
	jwt.StandardClaims
	UserID int
}

func GetPost(ctx *gin.Context) {
	postIDParam := ctx.Param("postid")
	resp, _ := get(postsServiceAddr + "/posts/" + postIDParam)
	ctx.JSON(resp.StatusCode, resp)
}

func GetPosts(ctx *gin.Context) {
	resp, _ := get(fmt.Sprintf("%s/posts?from=%v&total=%v&user=%v&sort=%v", 
	postsServiceAddr, 0, 10, "", ""))
	ctx.JSON(resp.StatusCode, resp)
}

func CreatePost(ctx *gin.Context) {
	defer ctx.Response.Body.Close()
	resp, _ := postProxy(postsServiceAddr + "/post/create", ctx.Response.Body)
	ctx.JSON(resp.StatusCode, resp)
}

func DeletePost(ctx *gin.Context) {
	userID, err := getUserID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	
	resp, _ := post(postsServiceAddr + "/post/delete", m)
	ctx.JSON(resp.StatusCode, resp)
}

func LikePost(ctx *gin.Context) {
	userID, err := getUserID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	postIDParam := ctx.Param("postid")
	resp, _ := post(postsServiceAddr + "/post/addlike/" + postIDParam, m)
	ctx.JSON(resp.StatusCode, resp)
}

func UnlikePost(ctx *gin.Context) {
	userID, err := getUserID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	postIDParam := ctx.Param("postid")
	resp, _ := post(postsServiceAddr + "/post/removelike/" + postIDParam, m)
	ctx.JSON(resp.StatusCode, resp)
}


func AddComment(ctx *gin.Context) {
	userID, err := getUserID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	postIDParam := ctx.Param("postid")
	resp, _ := post(postsServiceAddr + "/comment/create", m)
	ctx.JSON(resp.StatusCode, resp)
}

func LikeComment(ctx *gin.Context) {
	userID, err := getUserID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	postIDParam := ctx.Param("postid")
	resp, _ := post(postsServiceAddr + "/comment/addlike/" + postIDParam, m)
	ctx.JSON(resp.StatusCode, resp)
}

func UnlikeComment(ctx *gin.Context) {
	userID, err := getUserID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	postIDParam := ctx.Param("postid")
	resp, _ := post(postsServiceAddr + "/post/unlike/" + postIDParam, m)
	ctx.JSON(resp.StatusCode, resp)
}

func ClearComment(ctx *gin.Context) {
	userID, err := getUserID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	m, err := bindJSONAndAdd(ctx, gin.H{"UserID": userID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	postIDParam := ctx.Param("postid")
	resp, _ := post(postsServiceAddr + "/comment/clear/" + postIDParam, m)
	ctx.JSON(resp.StatusCode, resp)
}
func Search(ctx *gin.Context) {
	term := ctx.Query("term")
	page := ctx.Query("page")
	resp, _ := get(fmt.Sprintf("%s/search?from=%v&total=%v&term=%v", 
	searchServiceAddr, (page * 10) - 10, 10, term))
	ctx.JSON(resp.StatusCode, resp)
}

func RegisterGET(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID != 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func RegisterPOST(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID != 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	defer ctx.Response.Body.Close()
	resp, _ := postProxy(userServiceAddr + "/account", ctx.Response.Body)
	//add jwt
	ctx.JSON(resp.StatusCode, resp)
}

func LoginGET(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID != 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func LoginPOST(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID != 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	defer ctx.Response.Body.Close()
	resp, _ := postProxy(userServiceAddr + "/login", ctx.Response.Body)
	//add jwt
	ctx.JSON(resp.StatusCode, resp)
}

func ChangePassword(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}

	userID := getUserID(ctx)
	if userID != 0 {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	defer ctx.Response.Body.Close()
	resp, _ := postProxy(userServiceAddr + "/password", ctx.Response.Body)
	ctx.JSON(resp.StatusCode, resp)
}

func GetProfile(ctx *gin.Context) {
	userIDParam := ctx.Param("userid")
	resp, _ := get(userServiceAddr + "/" + userIDParam))
	ctx.JSON(resp.StatusCode, resp)
}

func UpdateProfile(ctx *gin.Context) {
	userID := getUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}

	opts := UpdateProfileOptions{}

	bio := ctx.PostForm("bio")
	if bio != "" {
		opts.Bio = bio
	}
	fileHeader, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	} else {
		objectName, err := mediaService.Upload(fileHeader.Filename)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
		opts.Avatar = objectName
	}

	err := userService.UpdateProfile(userID, opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func generateJWT(userID int) (string, error) {
	claims := JWTClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHMAC, claims)
	return token.SignedString(jwtSecretKey)
}

//jwt must be stored in authorization bearer
//format: Authorization: Bearer <token>
func getUserID(ctx *gin.Context) (int, error) {
	authHeader := ctx.GetHeader("Authorization")
	splitTokens := strings.Split(authHeader, "Bearer ")
	if len(splitTokens) != 2 {
		return 0, fmt.Errorf("invalid authorization header format")
	}
	reqToken := strings.TrimSpace(splitTokens[1])
	token, err := jwt.ParseWithClaims(reqToken, jwt.JWTClaims, func(t *jwt.Toke) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signingmethod")
		}
		
		return []byte(jwtKey), nil
	})

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims(jwt.JWTClaims); ok {
		userID := claims["userID"]
		if userID < 1 {
			return 0, fmt.Errorf("invalid userid")
		}
		return userID, nil
	}

	//should not happen
	return 0, fmt.Errorf("wrong claims type")
}