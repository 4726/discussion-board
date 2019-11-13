package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v6"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func assertJSON(t testing.TB, obj interface{}) string {
	b, err := json.Marshal(obj)
	assert.NoError(t, err)
	return string(b)
}

func fillDBTestData(t testing.TB, api *RestAPI) [][]byte {
	msg1 := []byte("hello")
	msg2 := []byte("world")
	opts := minio.PutObjectOptions{ContentType: "text/plain"}
	_, err := api.mc.PutObject(bucketName, "1", bytes.NewReader(msg1), int64(len(msg1)), opts)
	assert.NoError(t, err)
	_, err = api.mc.PutObject(bucketName, "2", bytes.NewReader(msg2), int64(len(msg2)), opts)
	assert.NoError(t, err)
	return [][]byte{msg1, msg2}
}

func queryDBTest(t testing.TB, api *RestAPI) [][]byte {
	contents := [][]byte{}

	doneCh := make(chan struct{})
	defer close(doneCh)
	objectCh := api.mc.ListObjects(bucketName, "", true, doneCh)
	for objectInfo := range objectCh {
		assert.NoError(t, objectInfo.Err)
		object, err := api.mc.GetObject(bucketName, objectInfo.Key, minio.GetObjectOptions{})
		assert.NoError(t, err)
		content, err := ioutil.ReadAll(object)
		assert.NoError(t, err)
		contents = append(contents, content)
	}

	return contents
}

func getCleanAPIForTesting(t testing.TB) *RestAPI {
	cfg, err := ConfigFromJSON("config_test.json")
	assert.NoError(t, err)
	api, err := NewRestAPI(cfg)
	assert.NoError(t, err)

	doneCh := make(chan struct{})
	defer close(doneCh)
	objectCh := api.mc.ListObjects(bucketName, "", true, doneCh)
	for objectInfo := range objectCh {
		assert.NoError(t, objectInfo.Err)
		err := api.mc.RemoveObject(bucketName, objectInfo.Key)
		assert.NoError(t, err)
	}

	return api
}

func TestUploadInvalidForm(t *testing.T) {
	api := getCleanAPIForTesting(t)

	contents := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/upload", nil)
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid form"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	contentsAfter := queryDBTest(t, api)
	assert.Equal(t, contents, contentsAfter)
}

func TestUploadNoMedia(t *testing.T) {
	api := getCleanAPIForTesting(t)

	contents := fillDBTestData(t, api)

	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	err := writer.Close()
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/upload", buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	api.engine.ServeHTTP(w, req)

	expected := ErrorResponse{"invalid form"}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, assertJSON(t, expected), w.Body.String())

	contentsAfter := queryDBTest(t, api)
	assert.Equal(t, contents, contentsAfter)
}

func TestUpload(t *testing.T) {
	api := getCleanAPIForTesting(t)

	contents := fillDBTestData(t, api)

	file, err := ioutil.TempFile("temp", "")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	_, err = file.Write([]byte("message"))
	assert.NoError(t, err)
	err = file.Close()
	assert.NoError(t, err)

	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	part, err := writer.CreateFormFile("media", file.Name())
	assert.NoError(t, err)
	part.Write([]byte("message"))
	err = writer.Close()
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/upload", buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `{"Name":`)

	contentsAfter := queryDBTest(t, api)
	contents = append(contents, []byte("message"))
	assert.ElementsMatch(t, contentsAfter, contents)
}

func TestRemoveDoesNotExist(t *testing.T) {
	api := getCleanAPIForTesting(t)

	contents := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/remove/3", nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	contentsAfter := queryDBTest(t, api)
	assert.Equal(t, contents, contentsAfter)
}

func TestRemove(t *testing.T) {
	api := getCleanAPIForTesting(t)

	contents := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/remove/1", nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{}", w.Body.String())

	expectedContents := [][]byte{contents[1]}

	contentsAfter := queryDBTest(t, api)
	assert.Equal(t, expectedContents, contentsAfter)
}

func TestInfo(t *testing.T) {
	api := getCleanAPIForTesting(t)

	contents := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/info", nil)
	api.engine.ServeHTTP(w, req)

	expected := map[string]interface{}{"StoreAddress": fmt.Sprintf("%s/%s/", api.mc.EndpointURL().String(), bucketName)}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, assertJSON(t, expected), w.Body.String())

	contentsAfter := queryDBTest(t, api)
	assert.Equal(t, contents, contentsAfter)
}

//tests if the bucket policy was set up correctly
func TestPublicReadBucket(t *testing.T) {
	api := getCleanAPIForTesting(t)

	contents := fillDBTestData(t, api)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/info", nil)
	api.engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	info := map[string]interface{}{}
	err := json.Unmarshal(w.Body.Bytes(), &info)
	assert.NoError(t, err)

	resp, err := http.Get(info["StoreAddress"].(string) + "1")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "text/plain", resp.Header["Content-Type"][0])
	assert.Equal(t, "hello", string(body))

	contentsAfter := queryDBTest(t, api)
	assert.Equal(t, contents, contentsAfter)
}
