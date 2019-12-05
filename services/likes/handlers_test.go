package main

import (
	"context"
	"fmt"
	"github.com/4726/discussion-board/services/likes/pb"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"os"
	"testing"
	"time"
)

var testGRPCApi *GRPCApi
var testAddr string

func TestMain(m *testing.M) {
	cfg, err := ConfigFromJSON("config_test.json")
	if err != nil {
		panic(err)
	}
	api, err := NewGRPCApi(cfg)
	if err != nil {
		panic(err)
	}
	testGRPCApi = api
	addr := fmt.Sprintf(":%v", cfg.ListenPort)
	testAddr = addr
	go api.Run(addr)
	time.Sleep(time.Second * 3)

	i := m.Run()
	//close server
	os.Exit(i)
}

func testSetup(t testing.TB) (pb.LikesClient, []PostLike, []CommentLike) {
	conn, err := grpc.Dial(testAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	// defer conn.Close()
	c := pb.NewLikesClient(conn)
	cleanDB(t)
	pLikes, cLikes := fillDBTestData(t)
	return c, pLikes, cLikes
}

func TestLikePostNoPostID(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{UserId: proto.Uint64(3)}
	_, err := c.LikePost(context.TODO(), req)
	assert.Error(t, err)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestLikePostNoUserID(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{Id: proto.Uint64(3)}
	_, err := c.LikePost(context.TODO(), req)
	assert.Error(t, err)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestLikePost(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{Id: proto.Uint64(1), UserId: proto.Uint64(3)}
	resp, err := c.LikePost(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.Total{Total: proto.Uint64(3)}
	assert.Equal(t, expected, resp)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assert.WithinDuration(t, pLikesAfter[2].CreatedAt, time.Now(), time.Second*10)
	pLikesAfter[2].CreatedAt = time.Time{}
	expectedPLikes := append(pLikes, PostLike{req.GetId(), req.GetUserId(), time.Time{}})
	assertPostsLikesEqual(t, expectedPLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestUnlikePostNoId(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{UserId: proto.Uint64(1)}
	_, err := c.UnlikePost(context.TODO(), req)
	assert.Error(t, err)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestUnlikePostNoUserId(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{Id: proto.Uint64(3)}
	_, err := c.UnlikePost(context.TODO(), req)
	assert.Error(t, err)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestUnlikePostDoesNotExist(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{Id: proto.Uint64(1), UserId: proto.Uint64(3)}
	resp, err := c.UnlikePost(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.Total{Total: proto.Uint64(2)}
	assert.Equal(t, expected, resp)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestUnlikePost(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{Id: proto.Uint64(1), UserId: proto.Uint64(1)}
	resp, err := c.UnlikePost(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.Total{Total: proto.Uint64(1)}
	assert.Equal(t, expected, resp)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes[1:], pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestLikeCommentNoCommentID(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{UserId: proto.Uint64(3)}
	_, err := c.LikeComment(context.TODO(), req)
	assert.Error(t, err)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestLikeCommentNoUserID(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{Id: proto.Uint64(3)}
	_, err := c.LikeComment(context.TODO(), req)
	assert.Error(t, err)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestLikeComment(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{Id: proto.Uint64(1), UserId: proto.Uint64(3)}
	resp, err := c.LikeComment(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.Total{Total: proto.Uint64(3)}
	assert.Equal(t, expected, resp)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assert.WithinDuration(t, cLikesAfter[2].CreatedAt, time.Now(), time.Second*10)
	cLikesAfter[2].CreatedAt = time.Time{}
	newCLike := CommentLike{req.GetId(), req.GetUserId(), time.Time{}}
	assertCommentsLikesEqual(t, []CommentLike{cLikes[0], cLikes[1], newCLike, cLikes[2]}, cLikesAfter)
}

func TestUnlikeCommentNoId(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{UserId: proto.Uint64(3)}
	_, err := c.UnlikeComment(context.TODO(), req)
	assert.Error(t, err)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestUnlikeCommentNoUserId(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{Id: proto.Uint64(3)}
	_, err := c.UnlikeComment(context.TODO(), req)
	assert.Error(t, err)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestUnlikeCommentDoesNotExist(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{Id: proto.Uint64(1), UserId: proto.Uint64(3)}
	resp, err := c.UnlikeComment(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.Total{Total: proto.Uint64(2)}
	assert.Equal(t, expected, resp)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestUnlikeComment(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDUserID{Id: proto.Uint64(1), UserId: proto.Uint64(2)}
	resp, err := c.UnlikeComment(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.Total{Total: proto.Uint64(1)}
	assert.Equal(t, expected, resp)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, []CommentLike{cLikes[0], cLikes[2]}, cLikesAfter)
}

func TestGetMultiplePostLikes(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	idLikes := []*pb.TotalLikes_IDLikes{}
	idLikes = append(idLikes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(1), Total: proto.Uint64(2)})
	idLikes = append(idLikes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(2), Total: proto.Uint64(0)})
	idLikes = append(idLikes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(3), Total: proto.Uint64(0)})
	expected := &pb.TotalLikes{IdLikes: idLikes}

	req := &pb.IDs{Id: []uint64{1, 2, 3}}
	resp, err := c.GetPostLikes(context.TODO(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestGetMultipleCommentLikes(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	idLikes := []*pb.TotalLikes_IDLikes{}
	idLikes = append(idLikes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(1), Total: proto.Uint64(2)})
	idLikes = append(idLikes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(2), Total: proto.Uint64(0)})
	idLikes = append(idLikes, &pb.TotalLikes_IDLikes{Id: proto.Uint64(3), Total: proto.Uint64(1)})
	expected := &pb.TotalLikes{IdLikes: idLikes}

	req := &pb.IDs{Id: []uint64{1, 2, 3}}
	resp, err := c.GetCommentLikes(context.TODO(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestHasMultiplePostLikesNoUserID(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDsUserID{Id: []uint64{1}}
	_, err := c.PostsHaveLike(context.TODO(), req)
	assert.Error(t, err)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestHasMultiplePostLikes(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	haveLikes := []*pb.HaveLikes_HaveLike{}
	haveLikes = append(haveLikes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(1), HasLike: proto.Bool(true)})
	haveLikes = append(haveLikes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(2), HasLike: proto.Bool(false)})
	haveLikes = append(haveLikes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(3), HasLike: proto.Bool(false)})
	expected := &pb.HaveLikes{HaveLikes: haveLikes}

	req := &pb.IDsUserID{Id: []uint64{1, 2, 3}, UserId: proto.Uint64(1)}
	resp, err := c.PostsHaveLike(context.TODO(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestHasMultipleCommentLikesNoUserID(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	req := &pb.IDsUserID{Id: []uint64{1}}
	_, err := c.CommentsHaveLike(context.TODO(), req)
	assert.Error(t, err)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestHasMultipleCommentLikes(t *testing.T) {
	c, pLikes, cLikes := testSetup(t)

	haveLikes := []*pb.HaveLikes_HaveLike{}
	haveLikes = append(haveLikes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(1), HasLike: proto.Bool(true)})
	haveLikes = append(haveLikes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(2), HasLike: proto.Bool(false)})
	haveLikes = append(haveLikes, &pb.HaveLikes_HaveLike{Id: proto.Uint64(3), HasLike: proto.Bool(true)})
	expected := &pb.HaveLikes{HaveLikes: haveLikes}

	req := &pb.IDsUserID{Id: []uint64{1, 2, 3}, UserId: proto.Uint64(2)}
	resp, err := c.CommentsHaveLike(context.TODO(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)

	pLikesAfter, cLikesAfter := queryDBTest(t)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func assertPostLikesEqual(t testing.TB, expected, actual PostLike) {
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
	expected.CreatedAt, actual.CreatedAt = time.Time{}, time.Time{}
	assert.Equal(t, expected, actual)
}

func assertPostsLikesEqual(t testing.TB, expected, actual []PostLike) {
	assert.Len(t, actual, len(expected))
	for i, v := range expected {
		assertPostLikesEqual(t, v, actual[i])
	}
}

func assertCommentLikesEqual(t testing.TB, expected, actual CommentLike) {
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
	expected.CreatedAt, actual.CreatedAt = time.Time{}, time.Time{}
	assert.Equal(t, expected, actual)
}

func assertCommentsLikesEqual(t testing.TB, expected, actual []CommentLike) {
	assert.Len(t, actual, len(expected))
	for i, v := range expected {
		assertCommentLikesEqual(t, v, actual[i])
	}
}

func cleanDB(t testing.TB) {
	testGRPCApi.handlers.db.Exec("TRUNCATE post_likes;")
	testGRPCApi.handlers.db.Exec("TRUNCATE comment_likes;")
}

func fillDBTestData(t testing.TB) ([]PostLike, []CommentLike) {
	p1 := addPostLike(t, 1, 1)
	p2 := addPostLike(t, 1, 2)
	c1 := addCommentLike(t, 1, 1)
	c2 := addCommentLike(t, 1, 2)
	c3 := addCommentLike(t, 3, 2)

	return []PostLike{p1, p2}, []CommentLike{c1, c2, c3}
}

func addPostLike(t testing.TB, postID, userID uint64) PostLike {
	like := PostLike{postID, userID, time.Now()}

	err := testGRPCApi.handlers.db.Save(&like).Error
	assert.NoError(t, err)
	return like
}

func addCommentLike(t testing.TB, commentID, userID uint64) CommentLike {
	like := CommentLike{commentID, userID, time.Now()}

	err := testGRPCApi.handlers.db.Save(&like).Error
	assert.NoError(t, err)
	return like
}

func queryDBTest(t testing.TB) ([]PostLike, []CommentLike) {
	postLikes := []PostLike{}
	commentLikes := []CommentLike{}
	assert.NoError(t, testGRPCApi.handlers.db.Find(&postLikes).Error)
	assert.NoError(t, testGRPCApi.handlers.db.Find(&commentLikes).Error)
	return postLikes, commentLikes
}
