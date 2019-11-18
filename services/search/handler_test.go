package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

type IndexForm struct {
	Title, Body string
	Id          uint
	UserID      uint
	Timestamp   int64
	Likes       int
}

type UpdateLikesForm struct {
	Id    uint
	Likes int
}

type DeletePostForm struct {
	Id uint
}

type UpdateLastUpdateForm struct {
	Id         uint
	LastUpdate int64
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func assertJSON(t testing.TB, obj interface{}) string {
	b, err := json.Marshal(obj)
	assert.NoError(t, err)
	return string(b)
}

func queryESC(t testing.TB, api *RestAPI) []Post {
	api.esc.client.Refresh().Index(api.esc.indexName).Do(context.TODO())
	posts := []Post{}
	query := elastic.NewMatchAllQuery()
	searchResult, err := api.esc.client.Search().
		Index(api.esc.indexName).
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

func getCleanAPIForTesting(t testing.TB) *RestAPI {
	api, err := NewRestAPI("testing", "http://127.0.0.1:9200")
	assert.NoError(t, err)

	query := elastic.NewMatchAllQuery()
	_, err = api.esc.client.DeleteByQuery().
		Index(api.esc.indexName).
		Query(query).
		Do(context.TODO())
	assert.NoError(t, err)

	return api
}

func fillESTestData(t testing.TB, api *RestAPI) []Post {
	p1 := Post{"my first post", "hello world", 1, 1, time.Now().Unix(), 0}
	indexForTesting(t, api, p1)
	p2 := Post{"post @2 world", "body #2", 2, 2, time.Now().Unix() + 10, 0}
	indexForTesting(t, api, p2)
	p3 := Post{"title3", "hello WORLd", 3, 3, time.Now().Unix() + 20, 0}
	indexForTesting(t, api, p3)
	return []Post{p1, p2, p3}
}

func indexForTesting(t testing.TB, api *RestAPI, p Post) {
	_, err := api.esc.client.Index().
		Index(api.esc.indexName).
		Refresh("wait_for").
		BodyJson(p).
		Do(context.TODO())
	assert.NoError(t, err)
}

func testInvalidBody(t *testing.T, form interface{}, route string) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", route, buffer)
	api.engine.ServeHTTP(w, req)

	expected := InvalidJSONBodyResponse

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestIndexInvalidBody(t *testing.T) {
	form := IndexForm{"", "body", 10, 10, 0, 0}
	testInvalidBody(t, form, "/index")
	form = IndexForm{"title", "", 10, 10, 0, 0}
	testInvalidBody(t, form, "/index")
	form = IndexForm{"title", "body", 0, 10, 0, 0}
	testInvalidBody(t, form, "/index")
	form = IndexForm{"title", "body", 10, 0, 0, 0}
	testInvalidBody(t, form, "/index")
}

func TestIndex(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	form := IndexForm{"title", "body", 10, 10, time.Now().Unix() + 30, 1}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/index", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Len(t, postsAfter, 4)
	assert.Equal(t, posts[0], postsAfter[0])
	assert.Equal(t, posts[1], postsAfter[1])
	assert.Equal(t, posts[2], postsAfter[2])
	addedPost := Post{form.Title, form.Body, form.Id, int(form.UserID), form.Timestamp, form.Likes}
	assert.Equal(t, addedPost, postsAfter[3])
}

func TestSearchInvalidFromQuery(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?from=a", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid query"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestSearchInvalidFromQuery2(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?from=-2", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid query"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestSearchInvalidTotalQuery(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?from=0&total=a", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid query"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestSearchInvalidTotalQuery2(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?from=0&total=0", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid query"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestSearchInvalidTotalQuery3(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?from=0&total=-1", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid query"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestSearchInvalidTermQuery(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?term=&from=0&total=10", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid query"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestSearch(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?term=hello&from=0&total=10", nil)
	api.engine.ServeHTTP(w, req)

	expected := []uint{posts[0].Id, posts[2].Id}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestSearch2(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?term=world&from=0&total=10", nil)
	api.engine.ServeHTTP(w, req)

	expected := []uint{posts[0].Id, posts[1].Id, posts[2].Id}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestUpdateLikesInvalidBody(t *testing.T) {
	form := UpdateLikesForm{0, 100}
	testInvalidBody(t, form, "/update/likes")
}

func TestUpdateLikesDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	form := UpdateLikesForm{5, 100}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update/likes", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestUpdateLikes(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	form := UpdateLikesForm{posts[1].Id, 100}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update/likes", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	expectedUpdatedPost := posts[1]
	expectedUpdatedPost.Likes = 100

	postsAfter := queryESC(t, api)
	assert.Len(t, postsAfter, 3)
	assert.Equal(t, posts[0], postsAfter[0])
	assert.Equal(t, expectedUpdatedPost, postsAfter[1])
	assert.Equal(t, posts[2], postsAfter[2])
}

func TestDeletePostInvalidJSON(t *testing.T) {
	form := DeletePostForm{0}
	testInvalidBody(t, form, "/deletepost")
}

func TestDeletePostDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	form := DeletePostForm{5}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/deletepost", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestDeletePost(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	form := DeletePostForm{posts[1].Id}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/deletepost", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	expectedUpdatedPosts := []Post{posts[0], posts[2]}

	postsAfter := queryESC(t, api)
	assert.Equal(t, expectedUpdatedPosts, postsAfter)
}

func TestUpdateLastUpdateInvalidJSON(t *testing.T) {
	form := UpdateLastUpdateForm{0, 1}
	testInvalidBody(t, form, "/update/lastupdate")
}

func TestUpdateLastUpdateDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	form := UpdateLastUpdateForm{5, time.Now().Unix()}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update/lastupdate", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestUpdateLastUpdate(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	newLastUpdate := time.Now().Unix()
	form := UpdateLastUpdateForm{posts[1].Id, newLastUpdate}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update/lastupdate", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	expectedUpdatedPost := posts[1]
	expectedUpdatedPost.LastUpdate = newLastUpdate

	postsAfter := queryESC(t, api)
	assert.Len(t, postsAfter, 3)
	assert.Equal(t, posts[0], postsAfter[0])
	assert.Equal(t, expectedUpdatedPost, postsAfter[1])
	assert.Equal(t, posts[2], postsAfter[2])
}
