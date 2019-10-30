package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func assertJSON(t testing.TB, obj interface{}) string {
	b, err := json.Marshal(obj)
	assert.NoError(t, err)
	return string(b)
}

func createAccountTest(t testing.TB, username, password string) {
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
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestGetProfileInvalidParam(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile/a", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid userid param"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestGetProfileDoesNotExist(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/profile/1", nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{}", w.Body.String())
}

func TestGetProfile(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	createAccountTest(t, "username", "password")

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
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestValidLoginWrongBodyFormat(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	buffer := bytes.NewBuffer([]byte("1"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", buffer)
	api.engine.ServeHTTP(w, req)

	expected := InvalidJSONBodyResponse

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestValidLoginUsernameDoesNotExist(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	form := LoginForm{"username", "password"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid login"}

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestValidLoginWrongPassword(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	createAccountTest(t, "username", "password")

	form := LoginForm{"username", "passwor"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid login"}

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestValidLogin(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	createAccountTest(t, "username", "password")

	form := LoginForm{"username", "password"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", buffer)
	api.engine.ServeHTTP(w, req)

	expected := map[string]interface{}{"userid": 1}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestCreateAccountWrongBodyFormat(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	buffer := bytes.NewBuffer([]byte("1"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/account", buffer)
	api.engine.ServeHTTP(w, req)

	expected := InvalidJSONBodyResponse

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestCreateAccountInvalidUsername(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	form := CreateAccountForm{"us", "password"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/account", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid username"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestCreateAccountInvalidPassword(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	form := CreateAccountForm{"username", "passw"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/account", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid password"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestCreateAccount(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	form := CreateAccountForm{"username", "password"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/account", buffer)
	api.engine.ServeHTTP(w, req)

	expected := map[string]interface{}{"userid": 1}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	form = CreateAccountForm{"username2", "password2"}
	buffer = bytes.NewBuffer([]byte(assertJSON(t, form)))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/account", buffer)
	api.engine.ServeHTTP(w, req)

	expected = map[string]interface{}{"userid": 2}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestUpdateProfileWrongBodyFormat(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	buffer := bytes.NewBuffer([]byte("1"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/profile/update", buffer)
	api.engine.ServeHTTP(w, req)

	expected := InvalidJSONBodyResponse

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestUpdateProfileNone(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	createAccountTest(t, "username", "password")

	form := UpdateProfileForm{1, "", ""}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/profile/update", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())
}

func TestUpdateProfileBio(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	createAccountTest(t, "username", "password")

	form := UpdateProfileForm{1, "hello world", ""}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/profile/update", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())
}

func TestUpdateProfileAvatarID(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	createAccountTest(t, "username", "password")

	form := UpdateProfileForm{1, "", "a"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/profile/update", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())
}

func TestUpdateProfile(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	createAccountTest(t, "username", "password")

	form := UpdateProfileForm{1, "hello world", "a"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/profile/update", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())
}

func TestChangePasswordWrongBodyFormat(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	//create account

	buffer := bytes.NewBuffer([]byte("1"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/password", buffer)
	api.engine.ServeHTTP(w, req)

	expected := InvalidJSONBodyResponse

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestChangePasswordUserDoesNotExist(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	form := ChangePasswordForm{1, "password", "newpassword"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/password", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{}", w.Body.String())
}

func TestChangePasswordInvalidOld(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	createAccountTest(t, "username", "password")

	form := ChangePasswordForm{1, "passwor", "newpassword"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/password", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid old password"}

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestChangePasswordInvalidNew(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	createAccountTest(t, "username", "password")

	form := ChangePasswordForm{1, "password", "newpa"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/password", buffer)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid new password"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())
}

func TestChangePassword(t *testing.T) {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)
	api.db.Exec("TRUNCATE auths;")
	api.db.Exec("TRUNCATE profiles;")

	createAccountTest(t, "username", "password")

	form := ChangePasswordForm{1, "password", "newpassword"}
	buffer := bytes.NewBuffer([]byte(assertJSON(t, form)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/password", buffer)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())
}
