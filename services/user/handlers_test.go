package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type LoginForm struct {
	Username, Password string
}

type CreateAccountForm struct {
	Username, Password string
}

type UpdateProfileForm struct {
	UserID        int
	Bio, AvatarID string
}

type ChangePasswordForm struct {
	UserID           int
	OldPass, NewPass string
}

func assertJSON(t testing.TB, obj interface{}) string {
	b, err := json.Marshal(obj)
	assert.NoError(t, err)
	return string(b)
}

func createAccountForTesting(t testing.TB, username, password string) ([]Auth, []Profile) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)

	form := CreateAccountForm{username, password}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/account", buffer)
	api.engine.ServeHTTP(w, req)

	expected := map[string]interface{}{"userid": 1}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	auths, profiles := queryDBTest(t, api)
	assert.Len(t, auths, 1)
	assert.Len(t, profiles, 1)
	auth := auths[0]
	profile := profiles[0]
	assert.Equal(t, 1, auth.UserID)
	assert.Equal(t, "username", auth.Username)
	assert.NotEqual(t, "password", auth.Password)
	assert.Equal(t, auth.CreatedAt, auth.UpdatedAt)
	assert.WithinDuration(t, time.Now(), auth.CreatedAt, time.Minute)
	assert.Equal(t, initialBio, profile.Bio)
	assert.Equal(t, initialAvatarID, profile.AvatarID)
	assert.Equal(t, profile.UserID, auth.UserID)
	assert.Equal(t, profile.Username, auth.Username)

	return auths, profiles
}

func queryDBTest(t testing.TB, api *RestAPI) ([]Auth, []Profile) {
	auths := []Auth{}
	profiles := []Profile{}
	assert.NoError(t, api.db.Find(&auths).Error)
	assert.NoError(t, api.db.Find(&profiles).Error)
	return auths, profiles
}

func getCleanAPIForTesting(t testing.TB) *RestAPI {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	return api
}

func testInvalidBody(t *testing.T, form interface{}, route string) {
	api := getCleanAPIForTesting(t)

	auths, profiles := createAccountForTesting(t, "username", "password")

	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", route, buffer)
	api.engine.ServeHTTP(w, req)

	expected := InvalidJSONBodyResponse

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	authsAfter, profilesAfter := queryDBTest(t, api)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestGetProfileInvalidParam(t *testing.T) {
	api := getCleanAPIForTesting(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile/a", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid userid param"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	auths, profiles := queryDBTest(t, api)
	assert.Len(t, auths, 0)
	assert.Len(t, profiles, 0)
}

func TestGetProfileDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile/1", nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	auths, profiles := queryDBTest(t, api)
	assert.Len(t, auths, 0)
	assert.Len(t, profiles, 0)
}

func TestGetProfile(t *testing.T) {
	api := getCleanAPIForTesting(t)

	auths, profiles := createAccountForTesting(t, "username", "password")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile/1", nil)
	api.engine.ServeHTTP(w, req)

	expected := Profile{
		UserID:   1,
		Username: "username",
		Bio:      "",
		AvatarID: "",
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	authsAfter, profilesAfter := queryDBTest(t, api)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestValidLoginInvalidBody(t *testing.T) {
	form := LoginForm{"", "password"}
	testInvalidBody(t, form, "/login")
	form = LoginForm{"username", ""}
	testInvalidBody(t, form, "/login")
}

func TestValidLoginUsernameDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := LoginForm{"username", "password"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid login"}

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	auths, profiles := queryDBTest(t, api)
	assert.Len(t, auths, 0)
	assert.Len(t, profiles, 0)
}

func TestValidLoginWrongPassword(t *testing.T) {
	api := getCleanAPIForTesting(t)

	auths, profiles := createAccountForTesting(t, "username", "password")

	form := LoginForm{"username", "passwor"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid login"}

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	authsAfter, profilesAfter := queryDBTest(t, api)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestValidLogin(t *testing.T) {
	api := getCleanAPIForTesting(t)

	auths, profiles := createAccountForTesting(t, "username", "password")

	form := LoginForm{"username", "password"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", buffer)
	api.engine.ServeHTTP(w, req)

	expected := map[string]interface{}{"userid": 1}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	authsAfter, profilesAfter := queryDBTest(t, api)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestCreateAccountInvalidBody(t *testing.T) {
	form := CreateAccountForm{"", "password"}
	testInvalidBody(t, form, "/account")
	form = CreateAccountForm{"username", ""}
	testInvalidBody(t, form, "/account")
}

func TestCreateAccountInvalidUsername(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := CreateAccountForm{"us", "password"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/account", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid username"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	auths, profiles := queryDBTest(t, api)
	assert.Len(t, auths, 0)
	assert.Len(t, profiles, 0)

}

func TestCreateAccountInvalidPassword(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := CreateAccountForm{"username", "passw"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/account", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid password"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	auths, profiles := queryDBTest(t, api)
	assert.Len(t, auths, 0)
	assert.Len(t, profiles, 0)
}

func TestCreateAccount(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := CreateAccountForm{"username", "password"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/account", buffer)
	api.engine.ServeHTTP(w, req)

	expected := map[string]interface{}{"userid": 1}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	form = CreateAccountForm{"username2", "password2"}
	buffer = bytes.NewBuffer([]byte(assertJSON(t, form)))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/account", buffer)
	api.engine.ServeHTTP(w, req)

	expected = map[string]interface{}{"userid": 2}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())
}

func TestUpdateProfileInvalidBody(t *testing.T) {
	form := UpdateProfileForm{0, "a", "a"}
	testInvalidBody(t, form, "/profile/update")
}

func TestUpdateProfileNone(t *testing.T) {
	api := getCleanAPIForTesting(t)

	auths, profiles := createAccountForTesting(t, "username", "password")

	form := UpdateProfileForm{1, "", ""}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/profile/update", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	authsAfter, profilesAfter := queryDBTest(t, api)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestUpdateProfileBio(t *testing.T) {
	api := getCleanAPIForTesting(t)

	auths, profiles := createAccountForTesting(t, "username", "password")

	form := UpdateProfileForm{1, "hello world 啊", ""}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/profile/update", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	authsAfter, profilesAfter := queryDBTest(t, api)
	profileAfter := profilesAfter[0]
	profile := profiles[0]
	assert.Len(t, profilesAfter, len(profiles))
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profile.UserID, profileAfter.UserID)
	assert.Equal(t, profile.Username, profileAfter.Username)
	assert.Equal(t, profileAfter.Bio, "hello world 啊")
	assert.Equal(t, profile.AvatarID, profileAfter.AvatarID)
}

func TestUpdateProfileAvatarID(t *testing.T) {
	api := getCleanAPIForTesting(t)

	auths, profiles := createAccountForTesting(t, "username", "password")

	form := UpdateProfileForm{1, "", "a"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/profile/update", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	authsAfter, profilesAfter := queryDBTest(t, api)
	profileAfter := profilesAfter[0]
	profile := profiles[0]
	assert.Len(t, profilesAfter, len(profiles))
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profile.UserID, profileAfter.UserID)
	assert.Equal(t, profile.Username, profileAfter.Username)
	assert.Equal(t, profile.Bio, profileAfter.Bio)
	assert.Equal(t, profileAfter.AvatarID, "a")
}

func TestUpdateProfile(t *testing.T) {
	api := getCleanAPIForTesting(t)

	auths, profiles := createAccountForTesting(t, "username", "password")

	form := UpdateProfileForm{1, "hello world", "a"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/profile/update", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	authsAfter, profilesAfter := queryDBTest(t, api)
	profileAfter := profilesAfter[0]
	profile := profiles[0]
	assert.Len(t, profilesAfter, len(profiles))
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profile.UserID, profileAfter.UserID)
	assert.Equal(t, profile.Username, profileAfter.Username)
	assert.Equal(t, profileAfter.AvatarID, "a")
	assert.Equal(t, profileAfter.Bio, "hello world")
}

func TestChangePasswordInvalidBody(t *testing.T) {
	form := ChangePasswordForm{0, "password", "newpassword"}
	testInvalidBody(t, form, "/password")
	form = ChangePasswordForm{1, "", "newpassword"}
	testInvalidBody(t, form, "/password")
	form = ChangePasswordForm{1, "password", ""}
	testInvalidBody(t, form, "/password")
}

func TestChangePasswordUserDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	form := ChangePasswordForm{1, "password", "newpassword"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/password", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	auths, profiles := queryDBTest(t, api)
	assert.Len(t, auths, 0)
	assert.Len(t, profiles, 0)
}

func TestChangePasswordInvalidOld(t *testing.T) {
	api := getCleanAPIForTesting(t)

	auths, profiles := createAccountForTesting(t, "username", "password")

	form := ChangePasswordForm{1, "passwor", "newpassword"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/password", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid old password"}

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	authsAfter, profilesAfter := queryDBTest(t, api)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestChangePasswordInvalidNew(t *testing.T) {
	api := getCleanAPIForTesting(t)

	auths, profiles := createAccountForTesting(t, "username", "password")

	form := ChangePasswordForm{1, "password", "newpa"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/password", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid new password"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	authsAfter, profilesAfter := queryDBTest(t, api)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestChangePassword(t *testing.T) {
	api := getCleanAPIForTesting(t)

	auths, profiles := createAccountForTesting(t, "username", "password")

	form := ChangePasswordForm{1, "password", "newpassword"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	time.Sleep(time.Second) //to test if updated_at field changes

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/password", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	authsAfter, profilesAfter := queryDBTest(t, api)
	auth := auths[0]
	authAfter := authsAfter[0]
	assert.Len(t, authsAfter, len(auths))
	assert.NotEqual(t, auth.Password, authAfter.Password)
	assert.NotEqual(t, auth.UpdatedAt, authAfter.UpdatedAt)
	assert.WithinDuration(t, time.Now(), authAfter.UpdatedAt, time.Minute)
	assert.Equal(t, auth.UserID, authAfter.UserID)
	assert.Equal(t, auth.Username, authAfter.Username)
	assert.Equal(t, auth.CreatedAt, authAfter.CreatedAt)
	assert.Equal(t, profiles, profilesAfter)
}
