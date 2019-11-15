package main

import (
	"encoding/json"
	"github.com/4726/discussion-board/services/posts/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func assertJSON(t testing.TB, obj interface{}) string {
	b, err := json.Marshal(obj)
	assert.NoError(t, err)
	return string(b)
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

func queryDBTest(t testing.TB, api *RestAPI) ([]models.Post, []models.Comment) {
	posts := []models.Post{}
	comments := []models.Comment{}
	assert.NoError(t, api.db.Preload("Comments").Find(&posts).Error)
	assert.NoError(t, api.db.Find(&comments).Error)
	return posts, comments
}

func getCleanAPIForTesting(t testing.TB) *RestAPI {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("DELETE FROM comments;")
	api.db.Exec("DELETE FROM posts;")
	api.db.Exec("ALTER TABLE posts AUTO_INCREMENT = 1;")
	api.db.Exec("ALTER TABLE comments AUTO_INCREMENT = 1;")

	return api
}

func fillDBTestData(t testing.TB, api *RestAPI) []models.Post {
	post := addPostForTesting(t, api, 2, "title", "hello world", 0)
	time.Sleep(time.Second)
	post2 := addPostForTesting(t, api, 1, "title2", "hello world 2", 5)
	time.Sleep(time.Second)
	comment1 := addCommentForTesting(t, api, 3, "my comment", 0, post2.ID)
	comment2 := addCommentForTesting(t, api, 4, "another comment", 0, post2.ID)
	post2.Comments = []models.Comment{comment1, comment2}
	post3 := addPostForTesting(t, api, 1, "title3", "hello world 3", 0)
	return []models.Post{post, post2, post3}
}

func addPostForTesting(t testing.TB, api *RestAPI, userID uint, title, body string, likes int) models.Post {
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

	err := api.db.Save(&post).Error
	assert.NoError(t, err)

	return post
}

func addCommentForTesting(t testing.TB, api *RestAPI, userID uint, body string, likes int, postID uint) models.Comment {
	created := time.Now()
	comment := models.Comment{
		PostID:    postID,
		UserID:    userID,
		Body:      body,
		Likes:     likes,
		CreatedAt: created,
	}

	err := api.db.Save(&comment).Error
	assert.NoError(t, err)

	return comment
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

	posts := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/1", nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	actualPost := models.Post{}
	json.Unmarshal(w.Body.Bytes(), &actualPost)
	assertPostEqual(t, posts[0], actualPost)

	postsAfter, _ := queryDBTest(t, api)
	assertPostsEqual(t, posts, postsAfter)
}

func TestGetPostsInvalidTotal(t *testing.T) {
	api := getCleanAPIForTesting(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts?total=a", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid query"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestGetPostsInvalidFrom(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts?total=1&from=a", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid query"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	postsAfter, _ := queryDBTest(t, api)
	assertPostsEqual(t, posts, postsAfter)
}

func TestGetPostsNoTotalQuery(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts?from=10", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid query"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	postsAfter, _ := queryDBTest(t, api)
	assertPostsEqual(t, posts, postsAfter)
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

	posts := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts?total=1&from=10&userid=10", nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())

	postsAfter, _ := queryDBTest(t, api)
	assertPostsEqual(t, posts, postsAfter)
}

func TestGetPostsSorted(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/posts?total=10&from=0&sort=created_at", nil)
	api.engine.ServeHTTP(w, req)

	expected := posts

	assert.Equal(t, http.StatusOK, w.Code)
	actualPosts := []models.Post{}
	json.Unmarshal(w.Body.Bytes(), &actualPosts)
	assert.Len(t, actualPosts, 3)
	assertMiniPostsEqual(t, expected, actualPosts)

	actualPosts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 3)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, actualPosts)
}

func TestGetPostsUserSorted(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/posts?total=10&from=0&userid=1&sort=created_at", nil)
	api.engine.ServeHTTP(w, req)

	expected := []models.Post{posts[1], posts[2]}

	assert.Equal(t, http.StatusOK, w.Code)
	actualPosts := []models.Post{}
	json.Unmarshal(w.Body.Bytes(), &actualPosts)
	assert.Len(t, actualPosts, 2)
	assertMiniPostsEqual(t, expected, actualPosts)

	actualPosts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 3)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, actualPosts)
}

func TestGetPostsUnSorted(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/posts?total=10&from=0", nil)
	api.engine.ServeHTTP(w, req)

	expected := []models.Post{posts[2], posts[1], posts[0]}

	assert.Equal(t, http.StatusOK, w.Code)
	actualPosts := []models.Post{}
	json.Unmarshal(w.Body.Bytes(), &actualPosts)
	assert.Len(t, actualPosts, 3)
	assertMiniPostsEqual(t, expected, actualPosts)

	actualPosts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 3)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, actualPosts)
}

func TestGetPostsUserUnSorted(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/posts?total=10&from=0&userid=1", nil)
	api.engine.ServeHTTP(w, req)

	expected := []models.Post{posts[2], posts[1]}

	assert.Equal(t, http.StatusOK, w.Code)
	actualPosts := []models.Post{}
	json.Unmarshal(w.Body.Bytes(), &actualPosts)
	assert.Len(t, actualPosts, 2)
	assertMiniPostsEqual(t, expected, actualPosts)

	actualPosts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 3)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, actualPosts)
}

func TestGetPostsTotal(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/posts?total=2&from=0", nil)
	api.engine.ServeHTTP(w, req)

	expected := []models.Post{posts[2], posts[1]}

	assert.Equal(t, http.StatusOK, w.Code)
	actualPosts := []models.Post{}
	json.Unmarshal(w.Body.Bytes(), &actualPosts)
	assert.Len(t, actualPosts, 2)
	assertMiniPostsEqual(t, expected, actualPosts)

	actualPosts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 3)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, actualPosts)
}

func TestGetPostsFrom(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/posts?total=2&from=1", nil)
	api.engine.ServeHTTP(w, req)

	expected := []models.Post{posts[1], posts[0]}

	assert.Equal(t, http.StatusOK, w.Code)
	actualPosts := []models.Post{}
	json.Unmarshal(w.Body.Bytes(), &actualPosts)
	assert.Len(t, actualPosts, 2)
	assertMiniPostsEqual(t, expected, actualPosts)

	actualPosts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 3)
	assert.Len(t, comments, 2)
	assertPostsEqual(t, posts, actualPosts)
}
