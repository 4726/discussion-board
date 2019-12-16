package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	pb "github.com/4726/discussion-board/services/media/pb"
	"github.com/golang/protobuf/proto"
	"github.com/minio/minio-go/v6"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var testApi *Api
var testAddr string
var testCFG Config

func TestMain(m *testing.M) {
	cfg, err := ConfigFromFile("config_test.json")
	if err != nil {
		panic(err)
	}
	api, err := NewApi(cfg)
	if err != nil {
		panic(err)
	}
	testApi = api
	testCFG = cfg
	addr := fmt.Sprintf(":%v", cfg.ListenPort)
	testAddr = addr
	go api.Run(addr)
	time.Sleep(time.Second * 3)

	i := m.Run()
	//close server
	os.Exit(i)
}

func testSetup(t testing.TB) (pb.MediaClient, [][]byte) {
	creds, err := credentials.NewClientTLSFromFile(testCFG.TLSCert, testCFG.TLSServerName)
	assert.NoError(t, err)
	conn, err := grpc.Dial(testAddr, grpc.WithTransportCredentials(creds))
	assert.NoError(t, err)
	// defer conn.Close()
	c := pb.NewMediaClient(conn)
	cleanDB(t)
	contents := fillDBTestData(t)
	return c, contents
}

func TestUploadNoMedia(t *testing.T) {
	c, contents := testSetup(t)

	req := &pb.UploadRequest{}
	_, err := c.Upload(context.TODO(), req)
	assert.Error(t, err)

	contentsAfter := queryDBTest(t)
	assert.Equal(t, contents, contentsAfter)
}

func TestUpload(t *testing.T) {
	c, contents := testSetup(t)

	req := &pb.UploadRequest{Media: []byte("message")}
	resp, err := c.Upload(context.TODO(), req)
	assert.NoError(t, err)
	assert.NotEqual(t, "", resp.Name)

	contentsAfter := queryDBTest(t)
	contents = append(contents, []byte("message"))
	assert.ElementsMatch(t, contentsAfter, contents)
}

func TestRemoveDoesNotExist(t *testing.T) {
	c, contents := testSetup(t)

	req := &pb.Name{}
	_, err := c.Remove(context.TODO(), req)
	assert.Error(t, err)

	contentsAfter := queryDBTest(t)
	assert.Equal(t, contents, contentsAfter)
}

func TestRemove(t *testing.T) {
	c, contents := testSetup(t)

	req := &pb.Name{Name: proto.String("1")}
	resp, err := c.Remove(context.TODO(), req)
	assert.NoError(t, err)
	assert.Equal(t, &pb.RemoveResponse{}, resp)

	expectedContents := [][]byte{contents[1]}
	contentsAfter := queryDBTest(t)
	assert.Equal(t, expectedContents, contentsAfter)
}

func TestInfo(t *testing.T) {
	c, contents := testSetup(t)

	req := &pb.InfoRequest{}
	resp, err := c.Info(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.InfoResponse{StoreAddress: proto.String(fmt.Sprintf("%s/%s/", testApi.mc.EndpointURL().String(), bucketName))}
	assert.Equal(t, expected, resp)

	contentsAfter := queryDBTest(t)
	assert.Equal(t, contents, contentsAfter)
}

//tests if the bucket policy was set up correctly
func TestPublicReadBucket(t *testing.T) {
	c, contents := testSetup(t)

	req := &pb.InfoRequest{}
	resp, err := c.Info(context.TODO(), req)
	assert.NoError(t, err)
	storeAddress := resp.StoreAddress

	getResp, err := http.Get(*storeAddress + "1")
	assert.NoError(t, err)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusOK, getResp.StatusCode)
	body, err := ioutil.ReadAll(getResp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "text/plain", getResp.Header["Content-Type"][0])
	assert.Equal(t, "hello", string(body))

	contentsAfter := queryDBTest(t)
	assert.Equal(t, contents, contentsAfter)
}

func fillDBTestData(t testing.TB) [][]byte {
	msg1 := []byte("hello")
	msg2 := []byte("world")
	opts := minio.PutObjectOptions{ContentType: "text/plain"}
	_, err := testApi.mc.PutObject(bucketName, "1", bytes.NewReader(msg1), int64(len(msg1)), opts)
	assert.NoError(t, err)
	_, err = testApi.mc.PutObject(bucketName, "2", bytes.NewReader(msg2), int64(len(msg2)), opts)
	assert.NoError(t, err)
	return [][]byte{msg1, msg2}
}

func queryDBTest(t testing.TB) [][]byte {
	contents := [][]byte{}

	doneCh := make(chan struct{})
	defer close(doneCh)
	objectCh := testApi.mc.ListObjects(bucketName, "", true, doneCh)
	for objectInfo := range objectCh {
		assert.NoError(t, objectInfo.Err)
		object, err := testApi.mc.GetObject(bucketName, objectInfo.Key, minio.GetObjectOptions{})
		assert.NoError(t, err)
		content, err := ioutil.ReadAll(object)
		assert.NoError(t, err)
		contents = append(contents, content)
	}

	return contents
}

func cleanDB(t testing.TB) {
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectCh := testApi.mc.ListObjects(bucketName, "", true, doneCh)
	for objectInfo := range objectCh {
		assert.NoError(t, objectInfo.Err)
		err := testApi.mc.RemoveObject(bucketName, objectInfo.Key)
		assert.NoError(t, err)
	}
}
