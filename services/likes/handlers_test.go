package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
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

func addPostLike(t testing.TB, api *RestAPI, postID, userID uint) PostLike {
	like := PostLike{postID, userID, time.Now()}

	err := api.db.Save(&like).Error
	assert.NoError(t, err)
	return like
}

func addCommentLike(t testing.TB, api *RestAPI, commentID, userID uint) CommentLike {
	like := CommentLike{commentID, userID, time.Now()}

	err := api.db.Save(&like).Error
	assert.NoError(t, err)
	return like
}

func fillDBTestData(t testing.TB, api *RestAPI) ([]PostLike, []CommentLike) {
	p1 := addPostLike(t, api, 1, 1)
	p2 := addPostLike(t, api, 1, 2)
	c1 := addCommentLike(t, api, 1, 1)
	c2 := addCommentLike(t, api, 1, 2)
	c3 := addCommentLike(t, api, 3, 2)

	return []PostLike{p1, p2}, []CommentLike{c1, c2, c3}
}

func queryDBTest(t testing.TB, api *RestAPI) ([]PostLike, []CommentLike) {
	postLikes := []PostLike{}
	commentLikes := []CommentLike{}
	assert.NoError(t, api.db.Find(&postLikes).Error)
	assert.NoError(t, api.db.Find(&commentLikes).Error)
	return postLikes, commentLikes
}

func getCleanAPIForTesting(t testing.TB) *RestAPI {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE post_likes;")
	api.db.Exec("TRUNCATE comment_likes;")

	return api
}

func TestLikePostInvalidForm(t *testing.T) {
	api := getCleanAPIForTesting(t)

	pLikes, cLikes := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/like", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid form"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	pLikesAfter, cLikesAfter := queryDBTest(t, api)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestLikePostNoPostID(t *testing.T) {
	api := getCleanAPIForTesting(t)

	pLikes, cLikes := fillDBTestData(t, api)

	form := PostLikeForm{0, 3}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/like", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid form"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	pLikesAfter, cLikesAfter := queryDBTest(t, api)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestLikePostNoUserID(t *testing.T) {
	api := getCleanAPIForTesting(t)

	pLikes, cLikes := fillDBTestData(t, api)

	form := PostLikeForm{1, 0}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/like", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid form"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	pLikesAfter, cLikesAfter := queryDBTest(t, api)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestLikePost(t *testing.T) {
	api := getCleanAPIForTesting(t)

	pLikes, cLikes := fillDBTestData(t, api)

	form := PostLikeForm{1, 3}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/like", buffer)
	api.engine.ServeHTTP(w, req)

	expected := gin.H{"total": 3}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	pLikesAfter, cLikesAfter := queryDBTest(t, api)
	assert.WithinDuration(t, pLikesAfter[2].CreatedAt, time.Now(), time.Second*10)
	pLikesAfter[2].CreatedAt = time.Time{}
	expectedPLikes := append(pLikes, PostLike{form.PostID, form.UserID, time.Time{}})
	assertPostsLikesEqual(t, expectedPLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestLikeCommentInvalidForm(t *testing.T) {
	api := getCleanAPIForTesting(t)

	pLikes, cLikes := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comment/like", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid form"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	pLikesAfter, cLikesAfter := queryDBTest(t, api)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestLikeCommentNoCommentID(t *testing.T) {
	api := getCleanAPIForTesting(t)

	pLikes, cLikes := fillDBTestData(t, api)

	form := CommentLikeForm{0, 3}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comment/like", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid form"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	pLikesAfter, cLikesAfter := queryDBTest(t, api)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestLikeCommentNoUserID(t *testing.T) {
	api := getCleanAPIForTesting(t)

	pLikes, cLikes := fillDBTestData(t, api)

	form := CommentLikeForm{3, 0}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comment/like", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid form"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	pLikesAfter, cLikesAfter := queryDBTest(t, api)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestLikeComment(t *testing.T) {
	api := getCleanAPIForTesting(t)

	pLikes, cLikes := fillDBTestData(t, api)

	form := CommentLikeForm{1, 3}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comment/like", buffer)
	api.engine.ServeHTTP(w, req)

	expected := gin.H{"total": 3}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	pLikesAfter, cLikesAfter := queryDBTest(t, api)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assert.WithinDuration(t, cLikesAfter[3].CreatedAt, time.Now(), time.Second*10)
	cLikesAfter[3].CreatedAt = time.Time{}
	expectedCLikes := append(cLikes, CommentLike{form.CommentID, form.UserID, time.Time{}})
	assertCommentsLikesEqual(t, expectedCLikes, cLikesAfter)
}

func TestGetMultiplePostLikesInvalidForm(t *testing.T) {
	api := getCleanAPIForTesting(t)

	pLikes, cLikes := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/likes", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid form"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	pLikesAfter, cLikesAfter := queryDBTest(t, api)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestGetMultiplePostLikes(t *testing.T) {
	api := getCleanAPIForTesting(t)

	pLikes, cLikes := fillDBTestData(t, api)

	form := IDsForm{[]uint{1, 2, 3}}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/likes", buffer)
	api.engine.ServeHTTP(w, req)

	expected := []IDLikes{}
	expected = append(expected, IDLikes{1, 2})
	expected = append(expected, IDLikes{2, 0})
	expected = append(expected, IDLikes{3, 0})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	pLikesAfter, cLikesAfter := queryDBTest(t, api)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestGetMultipleCommentLikesInvalidForm(t *testing.T) {
	api := getCleanAPIForTesting(t)

	pLikes, cLikes := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/comment/likes", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid form"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	pLikesAfter, cLikesAfter := queryDBTest(t, api)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}

func TestGetMultipleCommentLikes(t *testing.T) {
	api := getCleanAPIForTesting(t)

	pLikes, cLikes := fillDBTestData(t, api)

	form := IDsForm{[]uint{2, 3}}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/comment/likes", buffer)
	api.engine.ServeHTTP(w, req)

	expected := []IDLikes{}
	expected = append(expected, IDLikes{2, 0})
	expected = append(expected, IDLikes{3, 1})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	pLikesAfter, cLikesAfter := queryDBTest(t, api)
	assertPostsLikesEqual(t, pLikes, pLikesAfter)
	assertCommentsLikesEqual(t, cLikes, cLikesAfter)
}
