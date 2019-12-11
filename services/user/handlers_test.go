package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	pb "github.com/4726/discussion-board/services/user/pb"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var testApi *Api
var testAddr string

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
	addr := fmt.Sprintf(":%v", cfg.ListenPort)
	testAddr = addr
	go api.Run(addr)
	time.Sleep(time.Second * 3)

	i := m.Run()
	//close server
	os.Exit(i)
}

func testSetup(t testing.TB) (pb.UserClient, []Auth, []Profile) {
	conn, err := grpc.Dial(testAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	// defer conn.Close()
	c := pb.NewUserClient(conn)
	cleanDB(t)
	auths, profiles := fillDBTestData(t)
	return c, auths, profiles
}

func TestGetProfileDoesNotExist(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.UserId{UserId: proto.Uint64(10)}
	_, err := c.GetProfile(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestGetProfile(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.UserId{UserId: proto.Uint64(1)}
	resp, err := c.GetProfile(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.Profile{
		UserId:   proto.Uint64(1),
		Username: proto.String("aaaaaa"),
		Bio:      proto.String(""),
		AvatarId: proto.String(""),
	}
	assert.Equal(t, resp, expected)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestLoginPasswordRequired(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.LoginCredentials{Username: proto.String("")}
	_, err := c.Login(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestLoginUsernameRequired(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.LoginCredentials{Password: proto.String("")}
	_, err := c.Login(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestValidLoginUsernameDoesNotExist(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.LoginCredentials{Username: proto.String("q"), Password: proto.String("")}
	_, err := c.Login(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestValidLoginWrongPassword(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.LoginCredentials{Username: proto.String("aaaaaa"), Password: proto.String("111111")}
	_, err := c.Login(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestValidLogin(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.LoginCredentials{Username: proto.String("aaaaaa"), Password: proto.String("12345678")}
	resp, err := c.Login(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.UserId{UserId: proto.Uint64(1)}
	assert.Equal(t, expected, resp)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestCreateAccountPasswordRequired(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.LoginCredentials{Username: proto.String("")}
	_, err := c.CreateAccount(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestCreateAccountUsernameRequired(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.LoginCredentials{Password: proto.String("")}
	_, err := c.CreateAccount(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestCreateAccountInvalidUsername(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.LoginCredentials{Username: proto.String("as"), Password: proto.String("1111111")}
	_, err := c.CreateAccount(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestCreateAccountInvalidPassword(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.LoginCredentials{Username: proto.String("username"), Password: proto.String("passw")}
	_, err := c.CreateAccount(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestCreateAccountUsernameExists(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.LoginCredentials{Username: proto.String(auths[0].Username), Password: proto.String("password")}
	_, err := c.CreateAccount(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestCreateAccount(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.LoginCredentials{Username: proto.String("username"), Password: proto.String("password")}
	resp, err := c.CreateAccount(context.TODO(), req)
	assert.NoError(t, err)
	expected := &pb.UserId{UserId: proto.Uint64(4)}
	assert.Equal(t, expected, resp)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter[:3])
	assert.Equal(t, profiles, profilesAfter[:3])
	assert.Equal(t, uint64(4), authsAfter[3].UserID)
	assert.Equal(t, "username", authsAfter[3].Username)
	assert.Equal(t, authsAfter[3].CreatedAt, authsAfter[3].UpdatedAt)
	assert.WithinDuration(t, authsAfter[3].CreatedAt, time.Now(), time.Second*10)
	assert.Equal(t, profilesAfter[3].UserID, authsAfter[3].UserID)
	assert.Equal(t, profilesAfter[3].Username, authsAfter[3].Username)
	assert.Equal(t, initialBio, profilesAfter[3].Bio)
	assert.Equal(t, initialAvatarID, profilesAfter[3].AvatarID)
}

func TestUpdateProfileUserIdRequired(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.UpdateProfileRequest{}
	_, err := c.UpdateProfile(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestUpdateProfileNone(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.UpdateProfileRequest{UserId: proto.Uint64(1)}
	_, err := c.UpdateProfile(context.TODO(), req)
	assert.NoError(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestUpdateProfileBio(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.UpdateProfileRequest{UserId: proto.Uint64(1), Bio: proto.String("hello world 啊")}
	_, err := c.UpdateProfile(context.TODO(), req)
	assert.NoError(t, err)

	profiles[0].Bio = "hello world 啊"

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestUpdateProfileAvatarID(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.UpdateProfileRequest{UserId: proto.Uint64(1), AvatarId: proto.String("a")}
	_, err := c.UpdateProfile(context.TODO(), req)
	assert.NoError(t, err)

	profiles[0].AvatarID = "a"

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestUpdateProfile(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.UpdateProfileRequest{UserId: proto.Uint64(1), Bio: proto.String("hello world"), AvatarId: proto.String("a")}
	_, err := c.UpdateProfile(context.TODO(), req)
	assert.NoError(t, err)

	profiles[0].Bio = "hello world"
	profiles[0].AvatarID = "a"

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestChangePasswordUserIdRequired(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.ChangePasswordRequest{OldPass: proto.String("111111"), NewPass: proto.String("222222")}
	_, err := c.ChangePassword(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestChangePasswordOldPassRequired(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.ChangePasswordRequest{UserId: proto.Uint64(1), NewPass: proto.String("222222")}
	_, err := c.ChangePassword(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestChangePasswordNewPassRequired(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.ChangePasswordRequest{UserId: proto.Uint64(1), OldPass: proto.String("222222")}
	_, err := c.ChangePassword(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestChangePasswordUserDoesNotExist(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.ChangePasswordRequest{UserId: proto.Uint64(10), OldPass: proto.String("222222"), NewPass: proto.String("111111")}
	_, err := c.ChangePassword(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestChangePasswordInvalidOld(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.ChangePasswordRequest{UserId: proto.Uint64(1), OldPass: proto.String("222222"), NewPass: proto.String("111111")}
	_, err := c.ChangePassword(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}
func TestChangePasswordInvalidNew(t *testing.T) {
	c, auths, profiles := testSetup(t)

	req := &pb.ChangePasswordRequest{UserId: proto.Uint64(1), OldPass: proto.String("12345678"), NewPass: proto.String("newpa")}
	_, err := c.ChangePassword(context.TODO(), req)
	assert.Error(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths, authsAfter)
	assert.Equal(t, profiles, profilesAfter)
}

func TestChangePassword(t *testing.T) {
	c, auths, profiles := testSetup(t)

	time.Sleep(time.Second)

	req := &pb.ChangePasswordRequest{UserId: proto.Uint64(1), OldPass: proto.String("12345678"), NewPass: proto.String("newpassword")}
	_, err := c.ChangePassword(context.TODO(), req)
	assert.NoError(t, err)

	authsAfter, profilesAfter := queryDBTest(t)
	assert.Equal(t, auths[1:], authsAfter[1:])
	assert.Equal(t, profiles, profilesAfter)
	assert.NotEqual(t, auths[0].UpdatedAt, authsAfter[0].UpdatedAt)
	assert.WithinDuration(t, authsAfter[0].UpdatedAt, time.Now(), time.Second*5)
	assert.NotEqual(t, auths[0].Password, authsAfter[0].Password)
}

func cleanDB(t testing.TB) {
	testApi.db.Exec("TRUNCATE auths;")
	testApi.db.Exec("TRUNCATE profiles;")
}

func fillDBTestData(t testing.TB) ([]Auth, []Profile) {
	conn, err := grpc.Dial(testAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()
	c := pb.NewUserClient(conn)
	cleanDB(t)

	req := &pb.LoginCredentials{Username: proto.String("aaaaaa"), Password: proto.String("12345678")}
	_, err = c.CreateAccount(context.TODO(), req)
	assert.NoError(t, err)
	req = &pb.LoginCredentials{Username: proto.String("bbbbbb"), Password: proto.String("12345678")}
	_, err = c.CreateAccount(context.TODO(), req)
	assert.NoError(t, err)
	req = &pb.LoginCredentials{Username: proto.String("cccccc"), Password: proto.String("12345678")}
	_, err = c.CreateAccount(context.TODO(), req)
	assert.NoError(t, err)

	return queryDBTest(t)
}

func queryDBTest(t testing.TB) ([]Auth, []Profile) {
	auths := []Auth{}
	profiles := []Profile{}
	assert.NoError(t, testApi.db.Find(&auths).Error)
	assert.NoError(t, testApi.db.Find(&profiles).Error)
	return auths, profiles
}
