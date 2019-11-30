package main

import (
	"context"
	"fmt"
	"github.com/4726/discussion-board/services/posts/models"
	"github.com/4726/discussion-board/services/posts/write/pb"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"os"
	"testing"
	"time"
)

var testApi *Api
var testAddr string

func TestMain(m *testing.M) {
	cfg, err := ConfigFromJSON("config_test.json")
	if err != nil {
		panic(err)
	}
	api, err := NewApi(cfg)
	if err != nil {
		panic(err)
	}
	testApi = api
	addr := fmt.Sprintf(":%v", cfg.ListenPort)
	testAddr = addr
	go api.Run(addr)
	time.Sleep(time.Second * 3)

	i := m.Run()
	//close server
	os.Exit(i)
}

func testSetup(t testing.TB) (pb.PostsWriteClient, []models.Post) {
	conn, err := grpc.Dial(testAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	// defer conn.Close()
	c := pb.NewPostsWriteClient(conn)
	cleanDB(t)
	posts := fillDBTestData(t)
	return c, posts
}

func TestCreatePostEmptyUser(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.PostRequest{Title: proto.String("title"), Body: proto.String("body")}
	_, err := c.CreatePost(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestCreatePostEmptyTitle(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.PostRequest{UserId: proto.Uint64(1), Body: proto.String("body")}
	_, err := c.CreatePost(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestCreatePostEmptyBody(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.PostRequest{UserId: proto.Uint64(1), Title: proto.String("title")}
	_, err := c.CreatePost(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestCreatePost(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.PostRequest{UserId: proto.Uint64(1), Title: proto.String("first post"), Body: proto.String("hello world")}
	resp, err := c.CreatePost(context.TODO(), req)
	assert.NoError(t, err)
	assert.Equal(t, uint64(4), resp.GetPostId())

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter[:3])
	assert.Len(t, commentsAfter, 2)

	addedPost := postsAfter[3]
	assert.Equal(t, uint64(4), addedPost.ID)
	assert.Equal(t, uint64(1), addedPost.UserID)
	assert.Equal(t, "first post", addedPost.Title)
	assert.Equal(t, "hello world", addedPost.Body)
	assert.Equal(t, int64(0), addedPost.Likes)
	assert.Len(t, addedPost.Comments, 0)
	assert.Equal(t, addedPost.UpdatedAt, addedPost.CreatedAt)
	assert.WithinDuration(t, addedPost.CreatedAt, time.Now(), time.Second*10)
}

func TestDeletePostDoesNotExist(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.DeletePostRequest{PostId: proto.Uint64(10)}
	_, err := c.DeletePost(context.TODO(), req)
	assert.NoError(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestDeletePostNoPostID(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.DeletePostRequest{}
	_, err := c.DeletePost(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestDeletePostWithComments(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.DeletePostRequest{PostId: proto.Uint64(2)}
	_, err := c.DeletePost(context.TODO(), req)
	assert.NoError(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, []models.Post{posts[0], posts[2]}, postsAfter)
	assert.Len(t, commentsAfter, 0)
}

func TestDeletePostWithWrongUser(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.DeletePostRequest{PostId: proto.Uint64(1), UserId: proto.Uint64(3)}
	_, err := c.DeletePost(context.TODO(), req)
	assert.NoError(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestDeletePostWithRightUser(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.DeletePostRequest{PostId: proto.Uint64(1), UserId: proto.Uint64(2)}
	_, err := c.DeletePost(context.TODO(), req)
	assert.NoError(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, []models.Post{posts[1], posts[2]}, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestDeletePost(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.DeletePostRequest{PostId: proto.Uint64(1)}
	_, err := c.DeletePost(context.TODO(), req)
	assert.NoError(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, []models.Post{posts[1], posts[2]}, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestSetPostLikesNoPostID(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SetLikes{Likes: proto.Int64(1)}
	_, err := c.SetPostLikes(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestUpdatePostLikesNoLikes(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SetLikes{Id: proto.Uint64(1)}
	_, err := c.SetPostLikes(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestUpdatePostLikesDoesNotExist(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SetLikes{Id: proto.Uint64(10), Likes: proto.Int64(10)}
	_, err := c.SetPostLikes(context.TODO(), req)
	assert.NoError(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestUpdatePostLikes(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SetLikes{Id: proto.Uint64(1), Likes: proto.Int64(1)}
	_, err := c.SetPostLikes(context.TODO(), req)
	assert.NoError(t, err)

	posts[0].Likes = 1

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestCreateCommentNoPostID(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.CommentRequest{UserId: proto.Uint64(1), Body: proto.String("body")}
	_, err := c.CreateComment(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestCreateCommentNoUser(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.CommentRequest{PostId: proto.Uint64(1), Body: proto.String("body")}
	_, err := c.CreateComment(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestCreateCommentNoBody(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.CommentRequest{UserId: proto.Uint64(1), PostId: proto.Uint64(1)}
	_, err := c.CreateComment(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestCreateCommentPostDoesNotExist(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.CommentRequest{UserId: proto.Uint64(1), PostId: proto.Uint64(10), Body: proto.String("body")}
	_, err := c.CreateComment(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestCreateComment(t *testing.T) {
	c, posts := testSetup(t)
	time.Sleep(time.Second * 3) //sleep to check if UpdatedAt field changes

	req := &pb.CommentRequest{UserId: proto.Uint64(1), PostId: proto.Uint64(2), ParentId: proto.Uint64(1), Body: proto.String("body")}
	_, err := c.CreateComment(context.TODO(), req)
	assert.NoError(t, err)

	newComment := models.Comment{
		ID:       3,
		PostID:   2,
		ParentID: 1,
		UserID:   1,
		Body:     "body",
		Likes:    0,
	}
	posts[1].Comments = append(posts[1].Comments, newComment)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostEqual(t, posts[0], postsAfter[0])
	assertPostEqual(t, posts[2], postsAfter[2])
	assert.Len(t, commentsAfter, 3)

	assert.WithinDuration(t, time.Now(), postsAfter[1].Comments[2].CreatedAt, time.Second*10)
	postsAfter[1].Comments[2].CreatedAt = time.Time{}
	if postsAfter[1].UpdatedAt.Sub(posts[1].UpdatedAt) < time.Second*3 {
		assert.Fail(t, "UpdatedAt field on post did not update")
	}
	posts[1].UpdatedAt, postsAfter[1].UpdatedAt = time.Time{}, time.Time{}
	assertPostEqual(t, posts[1], postsAfter[1])
}

func TestClearCommentNoCommentID(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.ClearCommentRequest{}
	_, err := c.ClearComment(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestClearCommentDoesNotExist(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.ClearCommentRequest{CommentId: proto.Uint64(10)}
	_, err := c.ClearComment(context.TODO(), req)
	assert.NoError(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestClearCommentWithWrongUser(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.ClearCommentRequest{CommentId: proto.Uint64(1), UserId: proto.Uint64(1)}
	_, err := c.ClearComment(context.TODO(), req)
	assert.NoError(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestClearCommenttWithRightUser(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.ClearCommentRequest{CommentId: proto.Uint64(1), UserId: proto.Uint64(3)}
	_, err := c.ClearComment(context.TODO(), req)
	assert.NoError(t, err)

	posts[1].Comments[0].Body = ""

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestClearComment(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.ClearCommentRequest{CommentId: proto.Uint64(1), UserId: proto.Uint64(0)}
	_, err := c.ClearComment(context.TODO(), req)
	assert.NoError(t, err)

	posts[1].Comments[0].Body = ""

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestUpdateCommentLikesNoCommentID(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SetLikes{Likes: proto.Int64(1)}
	_, err := c.SetCommentLikes(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestUpdateCommentLikesNoLikes(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SetLikes{Id: proto.Uint64(1)}
	_, err := c.SetCommentLikes(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestUpdateCommentLikesLikesDoesNotExist(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SetLikes{Id: proto.Uint64(10), Likes: proto.Int64(10)}
	_, err := c.SetCommentLikes(context.TODO(), req)
	assert.NoError(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestUpdateCommentLikes(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SetLikes{Id: proto.Uint64(1), Likes: proto.Int64(3)}
	_, err := c.SetCommentLikes(context.TODO(), req)
	assert.NoError(t, err)

	posts[1].Comments[0].Likes = 3

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func cleanDB(t testing.TB) {
	testApi.db.Exec("DELETE FROM comments;")
	testApi.db.Exec("DELETE FROM posts;")
	testApi.db.Exec("ALTER TABLE posts AUTO_INCREMENT = 1;")
	testApi.db.Exec("ALTER TABLE comments AUTO_INCREMENT = 1;")
}

func queryDBTest(t testing.TB) ([]models.Post, []models.Comment) {
	posts := []models.Post{}
	comments := []models.Comment{}
	assert.NoError(t, testApi.db.Preload("Comments").Find(&posts).Error)
	assert.NoError(t, testApi.db.Find(&comments).Error)
	return posts, comments
}

func fillDBTestData(t testing.TB) []models.Post {
	post := addPostForTesting(t, 2, "title", "hello world", 0)
	time.Sleep(time.Second)
	post2 := addPostForTesting(t, 1, "title2", "hello world 2", 5)
	time.Sleep(time.Second)
	comment1 := addCommentForTesting(t, 3, "my comment", 0, post2.ID)
	comment2 := addCommentForTesting(t, 4, "another comment", 0, post2.ID)
	post2.Comments = []models.Comment{comment1, comment2}
	post3 := addPostForTesting(t, 1, "title3", "hello world 3", 0)
	return []models.Post{post, post2, post3}
}

func addPostForTesting(t testing.TB, userID uint64, title, body string, likes int64) models.Post {
	created := time.Now()
	post := models.Post{
		UserID:    userID,
		Title:     title,
		Body:      body,
		Likes:     likes,
		CreatedAt: created,
		UpdatedAt: created,
		Comments:  []models.Comment{},
	}

	err := testApi.db.Save(&post).Error
	assert.NoError(t, err)

	return post
}

func addCommentForTesting(t testing.TB, userID uint64, body string, likes int64, postID uint64) models.Comment {
	created := time.Now()
	comment := models.Comment{
		PostID:    postID,
		UserID:    userID,
		Body:      body,
		Likes:     likes,
		CreatedAt: created,
	}

	err := testApi.db.Save(&comment).Error
	assert.NoError(t, err)

	return comment
}

func assertPostEqual(t testing.TB, expected, actual models.Post) {
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, time.Second)
	assert.Equal(t, len(expected.Comments), len(actual.Comments))
	for i, v := range expected.Comments {
		assert.WithinDuration(t, v.CreatedAt, actual.Comments[i].CreatedAt, time.Second*2)
		v.CreatedAt, actual.Comments[i].CreatedAt = time.Time{}, time.Time{}
		assert.Equal(t, v, actual.Comments[i])
	}
	expected.CreatedAt, expected.UpdatedAt = time.Time{}, time.Time{}
	actual.CreatedAt, actual.UpdatedAt = time.Time{}, time.Time{}
	expected.Comments, actual.Comments = []models.Comment{}, []models.Comment{}
	assert.Equal(t, expected, actual)
}

func assertPostsEqual(t testing.TB, expected, actual []models.Post) {
	for i, v := range expected {
		assertPostEqual(t, v, actual[i])
	}
}
