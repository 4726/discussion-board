package main

import (
	"context"
	"regexp"
	"time"

	pb "github.com/4726/discussion-board/services/user/pb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	initialBio, initialAvatarID = "", ""
)

type Handlers struct {
	db *gorm.DB
}

func (h *Handlers) GetProfile(ctx context.Context, in *pb.UserId) (*pb.Profile, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	profile := Profile{}
	if err := h.db.First(&profile, in.GetUserId()).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, "user does not exist")
		}
		return nil, err
	}

	return &pb.Profile{
		UserId:   proto.Uint64(profile.UserID),
		Username: proto.String(profile.Username),
		Bio:      proto.String(profile.Bio),
		AvatarId: proto.String(profile.AvatarID),
	}, nil
}

func (h *Handlers) Login(ctx context.Context, in *pb.LoginCredentials) (*pb.UserId, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	auth := Auth{}
	if err := h.db.Where("username = ?", in.GetUsername()).First(&auth).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.Unauthenticated, "invalid login")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(in.GetPassword())); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid login")
	}

	return &pb.UserId{UserId: proto.Uint64(auth.UserID)}, nil
}

func (h *Handlers) CreateAccount(ctx context.Context, in *pb.LoginCredentials) (*pb.UserId, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	match, err := validUsername(in.GetUsername())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !match {
		return nil, status.Error(codes.InvalidArgument, "invalid username")
	}

	match, err = validPassword(in.GetPassword())
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, status.Error(codes.InvalidArgument, "invalid password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	created := time.Now()
	auth := Auth{
		Username:  in.GetUsername(),
		Password:  string(hash),
		CreatedAt: created,
		UpdatedAt: created,
	}
	if err := tx.Save(&auth).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, err.Error())
	}

	profile := Profile{
		UserID:   auth.UserID,
		Username: auth.Username,
		Bio:      initialBio,
		AvatarID: initialAvatarID,
	}
	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserId{UserId: proto.Uint64(auth.UserID)}, nil
}

func (h *Handlers) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	updates := map[string]interface{}{}
	updates["Bio"] = in.GetBio()
	updates["AvatarID"] = in.GetAvatarId()

	profile := Profile{UserID: in.GetUserId()}
	if err := h.db.Model(&profile).Updates(updates).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, "user id not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateProfileResponse{}, nil
}

func (h *Handlers) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	auth := Auth{}
	if err := tx.First(&auth, in.GetUserId()).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(in.GetOldPass())); err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	match, err := validPassword(in.GetNewPass())
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !match {
		tx.Rollback()
		return nil, status.Error(codes.InvalidArgument, "invalid new password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.GetNewPass()), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, err.Error())
	}

	updates := map[string]interface{}{
		"password":   string(hash),
		"updated_at": time.Now(),
	}

	if err := tx.Model(&auth).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ChangePasswordResponse{}, nil
}

func (h *Handlers) Check(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "client cancelled")
	}

	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING.Enum()}, nil
}

func validUsername(s string) (bool, error) {
	match, err := regexp.MatchString("^[A-Za-z0-9]{3,64}$", s)
	if err != nil {
		return false, err
	}
	if !match {
		return false, nil
	}

	return true, nil
}

func validPassword(s string) (bool, error) {
	match, err := regexp.MatchString("^[A-Za-z0-9]{6,64}$", s)
	if err != nil {
		return false, err
	}
	if !match {
		return false, nil
	}

	return true, nil
}
