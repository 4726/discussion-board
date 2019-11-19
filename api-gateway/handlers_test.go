package main

//just tests route interaction with/without jwt
//no integration testing

//comment out api.engine.Use(log.RequestMiddleware()) in NewRestAPI()
//when running tests. not sure why it causes tests to fail

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type statusCodeOpts struct {
	WithoutJWTStatusCode int
	WithJWTStatusCode    int
	HasEmptyResponseBefore bool
	HasEmptyResponseAfter bool
}

func getNewRestAPI(t *testing.T) *RestAPI {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	return api
}

func testJWTNotRequired(t *testing.T, method, route string) {
	opts := statusCodeOpts{http.StatusInternalServerError, http.StatusInternalServerError, true, true}
	testJWTStatusCode(t, method, route, opts)
}

func testJWTRequired(t *testing.T, method, route string) {
	opts := statusCodeOpts{http.StatusUnauthorized, http.StatusInternalServerError, true, true}
	testJWTStatusCode(t, method, route, opts)
}

func testJWTStatusCode(t *testing.T, method, route string, opts statusCodeOpts) {
	withoutJWT := opts.WithoutJWTStatusCode
	withJWT := opts.WithJWTStatusCode
	api := getNewRestAPI(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, route, nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, withoutJWT, w.Code)
	if opts.HasEmptyResponseBefore {
		assert.Equal(t, "{}", w.Body.String())
	} else {
		assert.NotEqual(t, "{}", w.Body.String())
	}

	//do again with jwt
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(method, route, nil)
	jwt, err := generateJWT(123)
	assert.NoError(t, err)
	req.Header["Authorization"] = []string{"Bearer " + jwt}
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, withJWT, w.Code)
	if opts.HasEmptyResponseAfter {
		assert.Equal(t, "{}", w.Body.String())
	} else {
		assert.NotEqual(t, "{}", w.Body.String())
	}
}

func TestJWTNotRequired(t *testing.T) {
	testJWTNotRequired(t, "GET", "/post/1")
	testJWTNotRequired(t, "GET", "/posts/1")
	testJWTNotRequired(t, "GET", "/search?term=hello&page=1")
	testJWTNotRequired(t, "GET", "/profile/1")

}

func TestJWTRequired(t *testing.T) {
	testJWTRequired(t, "POST", "/post")
	testJWTRequired(t, "DELETE", "/post/1")
	testJWTRequired(t, "POST", "/post/like")
	testJWTRequired(t, "POST", "/post/unlike")
	testJWTRequired(t, "POST", "/comment")
	testJWTRequired(t, "POST", "/comment/like")
	testJWTRequired(t, "POST", "/comment/unlike")
	testJWTRequired(t, "POST", "/comment/clear")
	testJWTRequired(t, "POST", "/changepassword")
}

func TestJWTCustom(t *testing.T) {
	testJWTStatusCode(t, "GET", "/register", statusCodeOpts{http.StatusOK, http.StatusBadRequest, true, false})
	testJWTStatusCode(t, "POST", "/register", statusCodeOpts{http.StatusInternalServerError, http.StatusBadRequest, true, true})
	testJWTStatusCode(t, "GET", "/login", statusCodeOpts{http.StatusOK, http.StatusBadRequest, true, false})
	testJWTStatusCode(t, "POST", "/login", statusCodeOpts{http.StatusInternalServerError, http.StatusBadRequest, true, true})
	testJWTStatusCode(t, "POST", "/profile/update", statusCodeOpts{http.StatusUnauthorized, http.StatusBadRequest, true, true})
}
