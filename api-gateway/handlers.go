package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/4726/discussion-board/api-gateway/pb/likes"
	postsread "github.com/4726/discussion-board/api-gateway/pb/posts-read"
	postswrite "github.com/4726/discussion-board/api-gateway/pb/posts-write"
	"github.com/4726/discussion-board/api-gateway/pb/search"
	"github.com/4726/discussion-board/api-gateway/pb/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

var jwtSecretKey = []byte("todosecretkey")

type JWTClaims struct {
	jwt.StandardClaims
	UserID uint64
}

func GetPost(ctx *gin.Context, clients GRPCClients) {
	postIdParam := ctx.Param("postid")
	postId, err := strconv.ParseUint(postIdParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	req := postsread.Id{Id: proto.Uint64(postId)}
	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	post, err := clients.PostsRead.GetFullPost(grpcCtx, &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.JSON(http.StatusNotFound, gin.H{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	m := structs.Map(post)

	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, post)
		return
	}

	req2 := likes.IDsUserID{UserId: proto.Uint64(userId), Id: []uint64{post.GetId()}}
	grpcCtx2, cancel2 := DefaultGRPCContext()
	defer cancel2()
	resp, err := clients.Likes.PostsHaveLike(grpcCtx2, &req2)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	m["has_like"] = resp.GetHaveLikes()[0].GetHasLike()

	commentIds := []uint64{}
	for _, v := range post.GetComments() {
		commentIds = append(commentIds, v.GetId())
	}
	req3 := likes.IDsUserID{UserId: proto.Uint64(userId), Id: commentIds}
	grpcCtx3, cancel3 := DefaultGRPCContext()
	defer cancel3()
	resp, err = clients.Likes.CommentsHaveLike(grpcCtx3, &req3)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	for _, v := range resp.GetHaveLikes() {
		id := v.GetId()
		hasLike := v.GetHasLike()
		comments := m["comments"].([]map[string]interface{})
		for i, comment := range comments {
			if comment["id"].(uint64) == id {
				comments[i]["has_like"] = hasLike
			}
		}
	}

	ctx.JSON(http.StatusOK, m)
}

func GetPosts(ctx *gin.Context, clients GRPCClients) {
	query := struct {
		Page   uint64 `form:"page" binding:"required"`
		UserId uint64 `form:"userid"`
	}{}
	if err := ctx.BindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	req := postsread.GetPostsQuery{
		From:   proto.Uint64(query.Page*10 - 10),
		Total:  proto.Uint64(10),
		UserId: proto.Uint64(query.UserId),
		Sort:   proto.String(""),
	}
	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	posts, err := clients.PostsRead.GetPosts(grpcCtx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, posts)
}

func CreatePost(ctx *gin.Context, clients GRPCClients) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	req := postswrite.PostRequest{}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)
	req.UserId = proto.Uint64(userId)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.PostsWrite.CreatePost(grpcCtx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func DeletePost(ctx *gin.Context, clients GRPCClients) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	req := postswrite.DeletePostRequest{}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)
	req.UserId = proto.Uint64(userId)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.PostsWrite.DeletePost(grpcCtx, &req)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument || status.Code(err) == codes.NotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	//should also remove from likes service

	ctx.JSON(http.StatusOK, resp)
}

func LikePost(ctx *gin.Context, clients GRPCClients) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	req := likes.IDUserID{}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)
	req.UserId = proto.Uint64(userId)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.Likes.LikePost(grpcCtx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	//should not matter much if it fails since
	//it can try again when another user likes/unlikes
	go func() {
		req2 := postswrite.SetLikes{
			Id:    proto.Uint64(req.GetId()),
			Likes: proto.Int64(int64(resp.GetTotal())),
		}
		grpcCtx, cancel := DefaultGRPCContext()
		defer cancel()
		_, _ = clients.PostsWrite.SetPostLikes(grpcCtx, &req2)
	}()

	ctx.JSON(http.StatusOK, resp)
}

func UnlikePost(ctx *gin.Context, clients GRPCClients) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	req := likes.IDUserID{}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)
	req.UserId = proto.Uint64(userId)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.Likes.UnlikePost(grpcCtx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	go func() {
		req2 := postswrite.SetLikes{
			Id:    proto.Uint64(req.GetId()),
			Likes: proto.Int64(int64(resp.GetTotal())),
		}
		grpcCtx, cancel := DefaultGRPCContext()
		defer cancel()
		_, _ = clients.PostsWrite.SetPostLikes(grpcCtx, &req2)
	
	}()

	ctx.JSON(http.StatusOK, resp)
}

func AddComment(ctx *gin.Context, clients GRPCClients) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	req := postswrite.CommentRequest{}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)
	req.UserId = proto.Uint64(userId)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.PostsWrite.CreateComment(grpcCtx, &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func LikeComment(ctx *gin.Context, clients GRPCClients) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	req := likes.IDUserID{}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)
	req.UserId = proto.Uint64(userId)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.Likes.LikeComment(grpcCtx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	go func() {
		req2 := postswrite.SetLikes{
			Id:    proto.Uint64(req.GetId()),
			Likes: proto.Int64(int64(resp.GetTotal())),
		}
		grpcCtx, cancel := DefaultGRPCContext()
		defer cancel()
		_, _ = clients.PostsWrite.SetCommentLikes(grpcCtx, &req2)
	}()

	ctx.JSON(http.StatusOK, resp)
}

func UnlikeComment(ctx *gin.Context, clients GRPCClients) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	req := likes.IDUserID{}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)
	req.UserId = proto.Uint64(userId)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.Likes.UnlikeComment(grpcCtx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	go func() {
		req2 := postswrite.SetLikes{
			Id:    proto.Uint64(req.GetId()),
			Likes: proto.Int64(int64(resp.GetTotal())),
		}
		grpcCtx, cancel := DefaultGRPCContext()
		defer cancel()
		_, _ = clients.PostsWrite.SetCommentLikes(grpcCtx, &req2)
	}()

	ctx.JSON(http.StatusOK, resp)
}

func ClearComment(ctx *gin.Context, clients GRPCClients) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	req := postswrite.ClearCommentRequest{}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)
	req.UserId = proto.Uint64(userId)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.PostsWrite.ClearComment(grpcCtx, &req)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument || status.Code(err) == codes.NotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func Search(ctx *gin.Context, clients GRPCClients) {
	form := struct {
		Term string `form:"term" binding:"required"`
		Page uint64 `form:"page" binding:"required"`
	}{}
	if err := ctx.BindQuery(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	req := search.SearchQuery{
		Term:  proto.String(form.Term),
		Total: proto.Uint64(10),
		From:  proto.Uint64(form.Page*10 - 10),
	}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.Search.Search(grpcCtx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	req2 := postsread.Ids{
		Id: resp.GetId(),
	}
	grpcCtx2, cancel2 := DefaultGRPCContext()
	defer cancel2()
	resp2, err := clients.PostsRead.GetPostsById(grpcCtx2, &req2)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, resp2)
}

func RegisterGET(ctx *gin.Context) {
	if userID, err := getUserId(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"user_id": userID})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func RegisterPOST(ctx *gin.Context, clients GRPCClients) {
	if _, err := getUserId(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	req := user.LoginCredentials{}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.User.CreateAccount(grpcCtx, &req)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	jwt, err := generateJWT(resp.GetUserId())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	m := structs.Map(resp)
	m["jwt"] = jwt

	ctx.JSON(http.StatusOK, m)
}

func LoginGET(ctx *gin.Context) {
	if userID, err := getUserId(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"user_id": userID})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func LoginPOST(ctx *gin.Context, clients GRPCClients) {
	if _, err := getUserId(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	req := user.LoginCredentials{}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.User.Login(grpcCtx, &req)
	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			ctx.JSON(http.StatusUnauthorized, gin.H{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	jwt, err := generateJWT(resp.GetUserId())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	m := structs.Map(resp)
	m["jwt"] = jwt
	ctx.JSON(http.StatusOK, m)
}

func ChangePassword(ctx *gin.Context, clients GRPCClients) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	req := user.ChangePasswordRequest{}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)
	req.UserId = proto.Uint64(userId)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.User.ChangePassword(grpcCtx, &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			//user id does not have an account
			ctx.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		if status.Code(err) == codes.InvalidArgument {
			//new password invalid
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		if status.Code(err) == codes.Unauthenticated {
			//wrong old password
			ctx.JSON(http.StatusUnauthorized, gin.H{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func GetProfile(ctx *gin.Context, clients GRPCClients) {
	isMine := false
	userId, _ := getUserId(ctx)
	profileIdParam := ctx.Param("userid")
	profileId, err := strconv.ParseUint(profileIdParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	if profileId == userId {
		isMine = true
	}

	req := user.UserId{UserId: proto.Uint64(profileId)}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.User.GetProfile(grpcCtx, &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.JSON(http.StatusNotFound, gin.H{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	m := structs.Map(resp)
	m["is_mine"] = isMine
	ctx.JSON(http.StatusOK, m)
}

func UpdateProfile(ctx *gin.Context, clients GRPCClients) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	req := user.UpdateProfileRequest{}
	defer ctx.Request.Body.Close()
	jsonpb.Unmarshal(ctx.Request.Body, &req)
	req.UserId = proto.Uint64(userId)

	grpcCtx, cancel := DefaultGRPCContext()
	defer cancel()
	resp, err := clients.User.UpdateProfile(grpcCtx, &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			//user id does not have a profile
			ctx.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func UserIdGET(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
	}

	ctx.JSON(http.StatusOK, gin.H{"user_id": userId})
}

func generateJWT(userId uint64) (string, error) {
	claims := JWTClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

//jwt must be stored in authorization bearer
//format: Authorization: Bearer <token>
func getUserId(ctx *gin.Context) (uint64, error) {
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
		userId := claims.UserID
		if userId < 1 {
			return 0, fmt.Errorf("invalid userid")
		}
		return uint64(userId), nil
	}

	//should not happen
	return 0, fmt.Errorf("wrong claims type")
}