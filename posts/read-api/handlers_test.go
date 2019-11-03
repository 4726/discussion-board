package main

import (
	"encoding/json"
	"fmt"
	"github.com/4726/discussion-board/posts/models"
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

func assertPostsEqual(t testing.TB, expected, actual models.Post) {
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, time.Second)
	expected.CreatedAt, expected.UpdatedAt = time.Time{}, time.Time{}
	actual.CreatedAt, actual.UpdatedAt = time.Time{}, time.Time{}
	assert.Equal(t, expected, actual)
}

func queryDBTest(t testing.TB, api *RestAPI) ([]models.Post, []models.Comment) {
	posts := []models.Post{}
	comments := []models.Comment{}
	assert.NoError(t, api.db.Find(&posts).Error)
	assert.NoError(t, api.db.Find(&comments).Error)
	return posts, comments
}

func getCleanAPIForTesting(t testing.TB) *RestAPI {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE comments;")
	api.db.Exec("TRUNCATE posts;")

	return api
}

func addPostForTesting(t testing.TB, api *RestAPI, username, title, body string) models.Post {
	created := time.Now()
	post := models.Post{
		User:      username,
		Title:     title,
		Body:      body,
		Likes:     0,
		CreatedAt: created,
		UpdatedAt: created,
	}

	err := api.db.Save(&post).Error
	assert.NoError(t, err)

	return post
}

func TestGetFullPostInvalidParam(t *testing.T) {
	api := getCleanAPIForTesting(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/a", nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestGetFullPostDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/1", nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestGetFullPost(t *testing.T) {
	api := getCleanAPIForTesting(t)

	post := addPostForTesting(t, api, "name", "title", "hello world")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/post/%v", post.PostID), nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	actualPost := models.Post{}
	json.Unmarshal(w.Body.Bytes(), &actualPost)
	assertPostsEqual(t, post, actualPost)

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 1)
	assert.Len(t, comments, 0)
	assertPostsEqual(t, posts[0], post)
}

func TestGetPostsInvalidTotal(t *testing.T) {
	api := getCleanAPIForTesting(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts?total=a", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid total query"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestGetPostsInvalidFrom(t *testing.T) {
	api := getCleanAPIForTesting(t)

	post := addPostForTesting(t, api, "name", "title", "hello world")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts?total=1&from=a", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid from query"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 1)
	assert.Len(t, comments, 0)
	assertPostsEqual(t, posts[0], post)
}

func TestGetPostsNoPosts(t *testing.T) {
	api := getCleanAPIForTesting(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts?total=1&from=10", nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestGetPostsUserNoPosts(t *testing.T) {
	api := getCleanAPIForTesting(t)

	post := addPostForTesting(t, api, "name", "title", "hello world")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts?total=1&from=10&user=asd", nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 1)
	assert.Len(t, comments, 0)
	assertPostsEqual(t, posts[0], post)
}

func TestGetPostsSorted(t *testing.T) {

}

func TestGetPostsUserSorted(t *testing.T) {
	api := getCleanAPIForTesting(t)

	post := addPostForTesting(t, api, "name", "title", "hello world")
	post2 := addPostForTesting(t, api, "asd", "title2", "hello world 2")
	post3 := addPostForTesting(t, api, "asd", "title3", "hello world 3")

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/posts?total=10&from=0&user=asd&sort=updated_at", nil)
	api.engine.ServeHTTP(w, req)

	expected := []models.Post{post3, post2}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 3)
	assert.Len(t, comments, 0)
	assertPostsEqual(t, posts[0], post)
	assertPostsEqual(t, posts[1], post2)
	assertPostsEqual(t, posts[2], post3)
}

func TestGetPostsUnSorted(t *testing.T) {

}

func TestGetPostsUserUnSorted(t *testing.T) {
	api := getCleanAPIForTesting(t)

	post := addPostForTesting(t, api, "name", "title", "hello world")
	post2 := addPostForTesting(t, api, "asd", "title2", "hello world 2")
	post3 := addPostForTesting(t, api, "asd", "title3", "hello world 3")

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/posts?total=10&from=0&user=asd", nil)
	api.engine.ServeHTTP(w, req)

	expected := []models.Post{post2, post3}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 3)
	assert.Len(t, comments, 0)
	assertPostsEqual(t, posts[0], post)
	assertPostsEqual(t, posts[1], post2)
	assertPostsEqual(t, posts[2], post3)
}
