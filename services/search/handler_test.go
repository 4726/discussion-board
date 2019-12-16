package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/4726/discussion-board/services/search/pb"
	"github.com/golang/protobuf/proto"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

var testApi *Api
var testAddr string
var testCFG Config

func TestMain(m *testing.M) {
	cfg, err := ConfigFromFile("config_test.json")
	if err != nil {
		panic(err)
	}
	api, err := NewApi(cfg)
	if err != nil {
		panic(err)
	}
	testApi = api
	testCFG = cfg
	addr := fmt.Sprintf(":%v", cfg.ListenPort)
	testAddr = addr
	go api.Run(addr)
	time.Sleep(time.Second * 3)

	i := m.Run()
	//close server
	os.Exit(i)
}

func testSetup(t testing.TB) (pb.SearchClient, []Post) {
	creds, err := credentials.NewClientTLSFromFile(testCFG.TLSCert, testCFG.TLSServerName)
	assert.NoError(t, err)
	conn, err := grpc.Dial(testAddr, grpc.WithTransportCredentials(creds))
	assert.NoError(t, err)
	// defer conn.Close()
	c := pb.NewSearchClient(conn)
	cleanDB(t)
	posts := fillESTestData(t)
	return c, posts
}

func TestIndexRequired(t *testing.T) {
	req := &pb.Post{
		Body:      proto.String("body"),
		UserId:    proto.Uint64(10),
		Id:        proto.Uint64(10),
		Timestamp: proto.Int64(0),
		Likes:     proto.Int64(0),
	}
	testIndexRequired(t, req)
	req = &pb.Post{
		Title:     proto.String("title"),
		UserId:    proto.Uint64(10),
		Id:        proto.Uint64(10),
		Timestamp: proto.Int64(0),
		Likes:     proto.Int64(0),
	}
	testIndexRequired(t, req)
	req = &pb.Post{
		Title:     proto.String("title"),
		Body:      proto.String("body"),
		Id:        proto.Uint64(10),
		Timestamp: proto.Int64(0),
		Likes:     proto.Int64(0),
	}
	testIndexRequired(t, req)
	req = &pb.Post{
		Title:     proto.String("title"),
		Body:      proto.String("body"),
		Id:        proto.Uint64(10),
		Timestamp: proto.Int64(0),
		Likes:     proto.Int64(0),
	}
	testIndexRequired(t, req)
	req = &pb.Post{
		Title:     proto.String("title"),
		Body:      proto.String("body"),
		UserId:    proto.Uint64(10),
		Timestamp: proto.Int64(0),
		Likes:     proto.Int64(0),
	}
	testIndexRequired(t, req)
}

func TestIndex(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Post{
		Title:     proto.String("title"),
		Body:      proto.String("body"),
		UserId:    proto.Uint64(10),
		Id:        proto.Uint64(10),
		Timestamp: proto.Int64(time.Now().Unix() + 30),
		Likes:     proto.Int64(1),
	}
	_, err := c.Index(context.TODO(), req)
	assert.NoError(t, err)

	expectedPost := Post{
		req.GetTitle(),
		req.GetBody(),
		req.GetId(),
		req.GetUserId(),
		req.GetTimestamp(),
		req.GetLikes(),
	}
	posts = append(posts, expectedPost)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestSearchNoTotal(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SearchQuery{Term: proto.String("term")}
	_, err := c.Search(context.TODO(), req)
	assert.Error(t, err)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestSearchNoTerm(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SearchQuery{Total: proto.Uint64(10)}
	_, err := c.Search(context.TODO(), req)
	assert.Error(t, err)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestSearch(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SearchQuery{Term: proto.String("hello"), Total: proto.Uint64(10), From: proto.Uint64(0)}
	resp, err := c.Search(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.SearchResult{Id: []uint64{posts[0].Id, posts[2].Id}}
	assert.Equal(t, expected, resp)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestSearch2(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SearchQuery{Term: proto.String("world"), Total: proto.Uint64(10), From: proto.Uint64(0)}
	resp, err := c.Search(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.SearchResult{Id: []uint64{posts[0].Id, posts[1].Id, posts[2].Id}}
	assert.Equal(t, expected, resp)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestSearch3(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.SearchQuery{Term: proto.String("world"), Total: proto.Uint64(10), From: proto.Uint64(2)}
	resp, err := c.Search(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.SearchResult{Id: []uint64{posts[2].Id}}
	assert.Equal(t, expected, resp)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestSetLikesNoId(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Likes{Likes: proto.Int64(100)}
	_, err := c.SetLikes(context.TODO(), req)
	assert.Error(t, err)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestSetLikesNoLikes(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Likes{Id: proto.Uint64(1)}
	_, err := c.SetLikes(context.TODO(), req)
	assert.Error(t, err)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestUpdateLikesDoesNotExist(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Likes{Id: proto.Uint64(5), Likes: proto.Int64(100)}
	_, err := c.SetLikes(context.TODO(), req)
	assert.NoError(t, err)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestUpdateLikes(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Likes{Id: proto.Uint64(posts[1].Id), Likes: proto.Int64(100)}
	_, err := c.SetLikes(context.TODO(), req)
	assert.NoError(t, err)

	posts[1].Likes = 100

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestDeletePostNoId(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Id{}
	_, err := c.DeletePost(context.TODO(), req)
	assert.Error(t, err)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestDeletePostDoesNotExist(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Id{Id: proto.Uint64(5)}
	_, err := c.DeletePost(context.TODO(), req)
	assert.NoError(t, err)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestDeletePost(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Id{Id: proto.Uint64(posts[1].Id)}
	_, err := c.DeletePost(context.TODO(), req)
	assert.NoError(t, err)

	posts = []Post{posts[0], posts[2]}

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestSetTimestampNoId(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Timestamp{Timestamp: proto.Int64(1)}
	_, err := c.SetTimestamp(context.TODO(), req)
	assert.Error(t, err)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestSetTimestampNoTimestamp(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Timestamp{Id: proto.Uint64(1)}
	_, err := c.SetTimestamp(context.TODO(), req)
	assert.Error(t, err)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestSetTimstampDoesNotExist(t *testing.T) {
	c, posts := testSetup(t)

	req := &pb.Timestamp{Id: proto.Uint64(5), Timestamp: proto.Int64(time.Now().Unix())}
	_, err := c.SetTimestamp(context.TODO(), req)
	assert.NoError(t, err)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func TestSetTimestamp(t *testing.T) {
	c, posts := testSetup(t)

	newTimestamp := time.Now().Unix()
	req := &pb.Timestamp{Id: proto.Uint64(posts[1].Id), Timestamp: proto.Int64(newTimestamp)}
	_, err := c.SetTimestamp(context.TODO(), req)
	assert.NoError(t, err)

	posts[1].LastUpdate = newTimestamp

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}

func queryESC(t testing.TB) []Post {
	testApi.esc.client.Refresh().Index(testApi.esc.indexName).Do(context.TODO())
	posts := []Post{}
	query := elastic.NewMatchAllQuery()
	searchResult, err := testApi.esc.client.Search().
		Index(testApi.esc.indexName).
		Query(query).
		SortBy(elastic.NewFieldSort("LastUpdate").Asc()).
		Do(context.TODO())
	assert.NoError(t, err)
	for _, hit := range searchResult.Hits.Hits {
		var p Post
		err := json.Unmarshal(hit.Source, &p)
		assert.NoError(t, err)
		posts = append(posts, p)
	}
	return posts
}

func cleanDB(t testing.TB) {
	query := elastic.NewMatchAllQuery()
	_, err := testApi.esc.client.DeleteByQuery().
		Index(testApi.esc.indexName).
		Query(query).
		Do(context.TODO())
	assert.NoError(t, err)
}

func fillESTestData(t testing.TB) []Post {
	p1 := Post{"my first post", "hello world", 1, 1, time.Now().Unix(), 0}
	indexForTesting(t, p1)
	p2 := Post{"post @2 world", "body #2", 2, 2, time.Now().Unix() + 10, 0}
	indexForTesting(t, p2)
	p3 := Post{"title3", "hello WORLd", 3, 3, time.Now().Unix() + 20, 0}
	indexForTesting(t, p3)
	return []Post{p1, p2, p3}
}

func indexForTesting(t testing.TB, p Post) {
	_, err := testApi.esc.client.Index().
		Index(testApi.esc.indexName).
		Refresh("wait_for").
		BodyJson(p).
		Do(context.TODO())
	assert.NoError(t, err)
}

func testIndexRequired(t *testing.T, req *pb.Post) {
	c, posts := testSetup(t)

	_, err := c.Index(context.TODO(), req)
	assert.Error(t, err)

	postsAfter := queryESC(t)
	assert.Equal(t, posts, postsAfter)
}
