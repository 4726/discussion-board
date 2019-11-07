package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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
	api, err := NewRestAPI("testing")
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
	p1 := Post{"my first post", "hello world", "user1", "id1", time.Now().Unix(), 0}
	indexForTesting(t, api, p1)
	p2 := Post{"post @2 world", "body #2", "user2", "id2", time.Now().Unix() + 10, 0}
	indexForTesting(t, api, p2)
	p3 := Post{"title3", "hello WORLd", "user3", "id3", time.Now().Unix() + 20, 0}
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

func TestIndexInvalidJSONForm(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	buffer := bytes.NewBuffer([]byte("1"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/index", buffer)
	api.engine.ServeHTTP(w, req)

	expected := InvalidJSONBodyResponse

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestIndex(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	form := IndexForm{"title", "body", "10", "10", time.Now().Unix() + 30, 1}
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
	addedPost := Post{form.Title, form.Body, form.User, form.Id, form.Timestamp, form.Likes}
	assert.Equal(t, addedPost, postsAfter[3])
}

func TestSearchInvalidFromQuery(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search?from=a", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid from query"}

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

	expected := ErrorResponse{"invalid total query"}

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

	expected := ErrorResponse{"invalid total query"}

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

	expected := ErrorResponse{"invalid term query"}

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

	expected := []string{posts[0].Id, posts[2].Id}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestUpdateLikesInvalidJSON(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	buffer := bytes.NewBuffer([]byte("1"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/updatelikes", buffer)
	api.engine.ServeHTTP(w, req)

	expected := InvalidJSONBodyResponse

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestUpdateLikesDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	form := UpdateLikesForm{"qweqwe", 100}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/updatelikes", buffer)
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
	req, _ := http.NewRequest("POST", "/updatelikes", buffer)
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
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	buffer := bytes.NewBuffer([]byte("1"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/deletepost", buffer)
	api.engine.ServeHTTP(w, req)

	expected := InvalidJSONBodyResponse

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	postsAfter := queryESC(t, api)
	assert.Equal(t, posts, postsAfter)
}

func TestDeletePostDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillESTestData(t, api)

	form := DeletePostForm{"qweqwe"}
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
