package main

import (
	"bytes"
	"encoding/json"
	"github.com/4726/discussion-board/services/posts/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type CreateForm struct {
	Title  string
	Body   string
	UserID uint
}

type DeleteForm struct {
	PostID uint
	UserID uint
}

type UpdateLikesForm struct {
	PostID uint
	Likes  int
}

type CreateCommentForm struct {
	PostID, ParentID uint
	UserID           uint
	Body             string
}

type ClearCommentForm struct {
	CommentID uint
	UserID    uint
}

type UpdateCommentLikesForm struct {
	CommentID uint
	Likes     int
}

func assertJSON(t testing.TB, obj interface{}) string {
	b, err := json.Marshal(obj)
	assert.NoError(t, err)
	return string(b)
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
	//cannot truncate because of foreign key constraints
	api.db.Exec("DELETE FROM comments;")
	api.db.Exec("DELETE FROM posts;")
	api.db.Exec("ALTER TABLE posts AUTO_INCREMENT = 1;")
	api.db.Exec("ALTER TABLE comments AUTO_INCREMENT = 1;")

	return api
}

func fillDBTestData(t testing.TB, api *RestAPI) []models.Post {
	post := addPostForTesting(t, api, 1, "title", "hello world", 0)
	time.Sleep(time.Second)
	post2 := addPostForTesting(t, api, 2, "title2", "hello world 2", 5)
	time.Sleep(time.Second)
	comment1 := addCommentForTesting(t, api, 3, "my comment", 0, post2.ID)
	comment2 := addCommentForTesting(t, api, 4, "another comment", 0, post2.ID)
	post2.Comments = []models.Comment{comment1, comment2}
	post3 := addPostForTesting(t, api, 2, "title3", "hello world 3", 0)
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

func TestCreatePostEmptyUser(t *testing.T) {
	form := CreateForm{"title", "body", 0}
	testInvalidBody(t, form, "/post/create")
}

func TestCreatePostEmptyTitle(t *testing.T) {
	form := CreateForm{"", "body", 1}
	testInvalidBody(t, form, "/post/create")
}

func TestCreatePostEmptyBody(t *testing.T) {
	form := CreateForm{"title", "", 1}
	testInvalidBody(t, form, "/post/create")
}

func TestCreatePost(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := CreateForm{"first post", "hello world", 1}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/create", buffer)
	api.engine.ServeHTTP(w, req)

	expected := map[string]interface{}{"postID": 1}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 1)
	assert.Len(t, comments, 0)
	post := posts[0]
	assert.Equal(t, uint(1), post.ID)
	assert.Equal(t, uint(1), post.UserID)
	assert.Equal(t, "first post", post.Title)
	assert.Equal(t, "hello world", post.Body)
	assert.Equal(t, 0, post.Likes)
	assert.Len(t, post.Comments, 0)
	assert.Equal(t, post.UpdatedAt, post.CreatedAt)
	assert.WithinDuration(t, post.CreatedAt, time.Now(), time.Second*10)
}

func TestDeletePostDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := DeleteForm{1, 0}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/delete", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestDeletePostNoPostID(t *testing.T) {
	form := DeleteForm{0, 0}
	testInvalidBody(t, form, "/post/delete")
}

func testInvalidBody(t *testing.T, form interface{}, route string) {
	api := getCleanAPIForTesting(t)

	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", route, buffer)
	api.engine.ServeHTTP(w, req)

	expected := InvalidJSONBodyResponse

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestDeletePostNoUser(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := DeleteForm{1, 0}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/delete", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestDeletePostWithComments(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	form := DeleteForm{2, 0}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/delete", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	postsAfter, comments := queryDBTest(t, api)
	assert.Len(t, postsAfter, 2)
	assert.Len(t, comments, 0)
	assertPostEqual(t, posts[0], postsAfter[0])
	assertPostEqual(t, posts[2], postsAfter[1])
}

func TestDeletePostWithWrongUser(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	form := DeleteForm{1, 2}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/delete", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	postsAfter, comments := queryDBTest(t, api)
	assert.Len(t, postsAfter, 3)
	assert.Len(t, comments, 2)
	assertPostEqual(t, posts[0], postsAfter[0])
	assertPostEqual(t, posts[1], postsAfter[1])
	assertPostEqual(t, posts[2], postsAfter[2])
}

func TestDeletePostWithRightUser(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	form := DeleteForm{1, 1}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/delete", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	postsAfter, comments := queryDBTest(t, api)
	assert.Len(t, postsAfter, 2)
	assert.Len(t, comments, 2)
	assertPostEqual(t, posts[1], postsAfter[0])
	assertPostEqual(t, posts[2], postsAfter[1])
}

func TestDeletePost(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	form := DeleteForm{1, 0}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/delete", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	postsAfter, comments := queryDBTest(t, api)
	assert.Len(t, postsAfter, 2)
	assert.Len(t, comments, 2)
	assertPostEqual(t, posts[1], postsAfter[0])
	assertPostEqual(t, posts[2], postsAfter[1])
}

func TestUpdatePostNoPostID(t *testing.T) {
	form := UpdateLikesForm{0, 1}
	testInvalidBody(t, form, "/post/likes")
}

func TestUpdatePostNoLikes(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := UpdateLikesForm{1, 0}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/likes", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestUpdatePostLikesDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := UpdateLikesForm{1, 1}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/likes", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestUpdatePostLikes(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	form := UpdateLikesForm{1, 1}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/likes", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	expectedUpdatedPost := posts[0]
	expectedUpdatedPost.Likes = 1

	postsAfter, comments := queryDBTest(t, api)
	assert.Len(t, postsAfter, 3)
	assert.Len(t, comments, 2)
	assertPostEqual(t, expectedUpdatedPost, postsAfter[0])
	assertPostEqual(t, posts[1], postsAfter[1])
	assertPostEqual(t, posts[2], postsAfter[2])
}

func TestCreateCommentNoPostID(t *testing.T) {
	form := CreateCommentForm{0, 0, 1, "body"}
	testInvalidBody(t, form, "/comment/create")
}

func TestCreateCommentNoUser(t *testing.T) {
	form := CreateCommentForm{1, 0, 0, "body"}
	testInvalidBody(t, form, "/comment/create")
}

func TestCreateCommentNoBody(t *testing.T) {
	form := CreateCommentForm{1, 0, 1, ""}
	testInvalidBody(t, form, "/comment/create")
}

func TestCreateCommentPostDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := CreateCommentForm{1, 0, 1, "body"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comment/create", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"post does not exist"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestCreateComment(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)
	time.Sleep(time.Second * 3) //sleep to check if UpdatedAt field changes

	form := CreateCommentForm{2, 1, 1, "body"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comment/create", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	expectedUpdatedPost := posts[1]
	expectedComment := models.Comment{
		ID:       3,
		PostID:   2,
		ParentID: 1,
		UserID:   1,
		Body:     "body",
		Likes:    0,
	}
	expectedUpdatedPost.Comments = append(expectedUpdatedPost.Comments, expectedComment)

	postsAfter, comments := queryDBTest(t, api)
	assert.Len(t, postsAfter, 3)
	assert.Len(t, comments, 3)
	assertPostEqual(t, posts[0], postsAfter[0])
	assertPostEqual(t, posts[2], postsAfter[2])

	assert.WithinDuration(t, time.Now(), postsAfter[1].Comments[2].CreatedAt, time.Second*10)
	postsAfter[1].Comments[2].CreatedAt = time.Time{}
	if postsAfter[1].UpdatedAt.Sub(expectedUpdatedPost.UpdatedAt) < time.Second*3 {
		assert.Fail(t, "UpdatedAt field on post did not update")
	}
	expectedUpdatedPost.UpdatedAt, postsAfter[1].UpdatedAt = time.Time{}, time.Time{}
	assertPostEqual(t, expectedUpdatedPost, postsAfter[1])
}

func TestClearCommentNoCommentID(t *testing.T) {
	form := ClearCommentForm{0, 0}
	testInvalidBody(t, form, "/comment/clear")
}

func TestClearCommentDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := ClearCommentForm{1, 0}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comment/clear", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestClearCommentWithWrongUser(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	form := ClearCommentForm{1, 1}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comment/clear", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	postsAfter, comments := queryDBTest(t, api)
	assert.Len(t, postsAfter, 3)
	assert.Len(t, comments, 2)
	assertPostEqual(t, posts[0], postsAfter[0])
	assertPostEqual(t, posts[1], postsAfter[1])
	assertPostEqual(t, posts[2], postsAfter[2])
}

func TestClearCommenttWithRightUser(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	form := ClearCommentForm{1, 3}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comment/clear", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	expectedUpdatedPost := posts[1]
	expectedUpdatedPost.Comments[0].Body = ""

	postsAfter, comments := queryDBTest(t, api)
	assert.Len(t, postsAfter, 3)
	assert.Len(t, comments, 2)
	assertPostEqual(t, posts[0], postsAfter[0])
	assertPostEqual(t, posts[2], postsAfter[2])
	assertPostEqual(t, expectedUpdatedPost, postsAfter[1])
}

func TestClearComment(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	form := ClearCommentForm{1, 0}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comment/clear", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	expectedUpdatedPost := posts[1]
	expectedUpdatedPost.Comments[0].Body = ""

	postsAfter, comments := queryDBTest(t, api)
	assert.Len(t, postsAfter, 3)
	assert.Len(t, comments, 2)
	assertPostEqual(t, posts[0], postsAfter[0])
	assertPostEqual(t, posts[2], postsAfter[2])
	assertPostEqual(t, expectedUpdatedPost, postsAfter[1])
}

func TestUpdateCommentNoCommentID(t *testing.T) {
	form := UpdateCommentLikesForm{0, 1}
	testInvalidBody(t, form, "/comment/likes")
}

func TestUpdateCommentLikesDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := UpdateCommentLikesForm{1, 1}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comment/likes", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	posts, comments := queryDBTest(t, api)
	assert.Len(t, posts, 0)
	assert.Len(t, comments, 0)
}

func TestUpdateCommentLikes(t *testing.T) {
	api := getCleanAPIForTesting(t)

	posts := fillDBTestData(t, api)

	form := UpdateCommentLikesForm{1, 3}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/comment/likes", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	expectedUpdatedPost := posts[1]
	expectedUpdatedPost.Comments[0].Likes = 3

	postsAfter, comments := queryDBTest(t, api)
	assert.Len(t, postsAfter, 3)
	assert.Len(t, comments, 2)
	assertPostEqual(t, posts[0], postsAfter[0])
	assertPostEqual(t, posts[2], postsAfter[2])
	assertPostEqual(t, expectedUpdatedPost, postsAfter[1])
}
