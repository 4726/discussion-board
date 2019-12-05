package main

import (
	"context"
	"fmt"
	"github.com/4726/discussion-board/services/posts/models"
	"github.com/4726/discussion-board/services/posts/read/pb"
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

func testSetup(t testing.TB) (pb.PostsReadClient, []models.Post) {
	conn, err := grpc.Dial(testAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	// defer conn.Close()
	c := pb.NewPostsReadClient(conn)
	cleanDB(t)
	posts := fillDBTestData(t)
	return c, posts
}

func TestGetFullPostNoId(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Id{}
	_, err := c.GetFullPost(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestGetFullPostDoesNotExist(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Id{Id: proto.Uint64(5)}
	_, err := c.GetFullPost(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestGetFullPost(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Id{Id: proto.Uint64(1)}
	resp, err := c.GetFullPost(context.TODO(), req)
	assert.NoError(t, err)
	assertPostEqual(t, posts[0], protoPostToModelPost(resp))

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestGetPostsNoTotal(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.GetPostsQuery{}
	_, err := c.GetPosts(context.TODO(), req)
	assert.Error(t, err)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestGetPostsNoPosts(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.GetPostsQuery{Total: proto.Uint64(1), From: proto.Uint64(10)}
	resp, err := c.GetPosts(context.TODO(), req)
	assert.NoError(t, err)
	assert.Len(t, resp.Posts, 0)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestGetPostsUserNoPosts(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.GetPostsQuery{Total: proto.Uint64(1), UserId: proto.Uint64(10)}
	resp, err := c.GetPosts(context.TODO(), req)
	assert.NoError(t, err)
	assert.Len(t, resp.Posts, 0)

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestGetPostsSorted(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.GetPostsQuery{Total: proto.Uint64(10), Sort: proto.String("created_at")}
	resp, err := c.GetPosts(context.TODO(), req)
	assert.NoError(t, err)
	expected := posts
	assertMiniPostsEqual(t, expected, protoPostsToModelPosts(resp.Posts))

	postsAfter, comments := queryDBTest(t)
	assert.Len(t, posts, 3)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, postsAfter)
}

func TestGetPostsUserSorted(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.GetPostsQuery{Total: proto.Uint64(10), Sort: proto.String("created_at"), UserId: proto.Uint64(1)}
	resp, err := c.GetPosts(context.TODO(), req)
	assert.NoError(t, err)
	expected := []models.Post{posts[1], posts[2]}
	assertMiniPostsEqual(t, expected, protoPostsToModelPosts(resp.Posts))

	psotsAfter, comments := queryDBTest(t)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, psotsAfter)
}

func TestGetPostsUnSorted(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.GetPostsQuery{Total: proto.Uint64(10)}
	resp, err := c.GetPosts(context.TODO(), req)
	assert.NoError(t, err)
	expected := []models.Post{posts[2], posts[1], posts[0]}
	assertMiniPostsEqual(t, expected, protoPostsToModelPosts(resp.Posts))

	postsAfter, comments := queryDBTest(t)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, postsAfter)
}

func TestGetPostsUserUnSorted(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.GetPostsQuery{Total: proto.Uint64(10), UserId: proto.Uint64(1)}
	resp, err := c.GetPosts(context.TODO(), req)
	assert.NoError(t, err)
	expected := []models.Post{posts[2], posts[1]}
	assertMiniPostsEqual(t, expected, protoPostsToModelPosts(resp.Posts))

	postsAfter, comments := queryDBTest(t)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, postsAfter)
}

func TestGetPostsTotal(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.GetPostsQuery{Total: proto.Uint64(2)}
	resp, err := c.GetPosts(context.TODO(), req)
	assert.NoError(t, err)
	expected := []models.Post{posts[2], posts[1]}
	assertMiniPostsEqual(t, expected, protoPostsToModelPosts(resp.Posts))

	postsAfter, comments := queryDBTest(t)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, postsAfter)
}

func TestGetPostsFrom(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.GetPostsQuery{Total: proto.Uint64(2), From: proto.Uint64(1)}
	resp, err := c.GetPosts(context.TODO(), req)
	assert.NoError(t, err)
	expected := []models.Post{posts[1], posts[0]}
	assertMiniPostsEqual(t, expected, protoPostsToModelPosts(resp.Posts))

	postsAfter, commentsAfter := queryDBTest(t)
	assertPostsEqual(t, posts, postsAfter)
	assert.Len(t, commentsAfter, 2)
}

func TestGetPostsByIdNone(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Ids{}
	resp, err := c.GetPostsById(context.TODO(), req)
	assert.NoError(t, err)
	assert.Len(t, resp.Posts, 0)

	postsAfter, comments := queryDBTest(t)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, postsAfter)
}

func TestGetMultiplePostsDoesNotExist(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Ids{Id: []uint64{5, 10}}
	resp, err := c.GetPostsById(context.TODO(), req)
	assert.NoError(t, err)
	assert.Len(t, resp.Posts, 0)

	postsAfter, comments := queryDBTest(t)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, postsAfter)
}

func TestGetMultiplePosts(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Ids{Id: []uint64{2, 3}}
	resp, err := c.GetPostsById(context.TODO(), req)
	assert.NoError(t, err)
	expected := []models.Post{posts[1], posts[2]}
	assertMiniPostsEqual(t, expected, protoPostsToModelPosts(resp.Posts))

	postsAfter, comments := queryDBTest(t)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, postsAfter)
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

func assertPostEqual(t testing.TB, expected, actual models.Post) {
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, time.Second)
	assert.Equal(t, len(expected.Comments), len(actual.Comments))
	for i, v := range expected.Comments {
		assert.WithinDuration(t, v.CreatedAt, actual.Comments[i].CreatedAt, time.Second)
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

func assertMiniPostEqual(t testing.TB, expected, actual models.Post) {
	expected.Body = ""
	actual.Body = ""
	expected.Comments = []models.Comment{}
	actual.Comments = []models.Comment{}
	assertPostEqual(t, expected, actual)
}

func assertMiniPostsEqual(t testing.TB, expected, actual []models.Post) {
	for i, v := range expected {
		assertMiniPostEqual(t, v, actual[i])
	}
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

func protoPostsToModelPosts(posts []*pb.Post) []models.Post {
	modelPosts := []models.Post{}

	for _, v := range posts {
		modelPosts = append(modelPosts, protoPostToModelPost(v))
	}

	return modelPosts
}

func protoPostToModelPost(post *pb.Post) models.Post {
	modelComments := []models.Comment{}

	for _, v := range post.Comments {
		modelComment := protoCommentToModelComment(v)
		modelComments = append(modelComments, modelComment)
	}

	return models.Post{
		ID:        *post.Id,
		UserID:    *post.UserId,
		Title:     *post.Title,
		Body:      *post.Body,
		Likes:     *post.Likes,
		CreatedAt: time.Unix(*post.CreatedAt, 0),
		UpdatedAt: time.Unix(*post.UpdatedAt, 0),
		Comments:  modelComments,
	}
}

func protoCommentToModelComment(comment *pb.Comment) models.Comment {
	return models.Comment{
		ID:        *comment.Id,
		PostID:    *comment.PostId,
		ParentID:  *comment.ParentId,
		UserID:    *comment.UserId,
		Body:      *comment.Body,
		CreatedAt: time.Unix(*comment.CreatedAt, 0),
		Likes:     *comment.Likes,
	}
}
