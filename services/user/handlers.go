package main

import (
	"context"
	"fmt"
	"github.com/4726/discussion-board/services/user/pb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

const (
	initialBio, initialAvatarID = "", ""
)

type Handlers struct {
	db *gorm.DB
}

func (h *Handlers) GetProfile(ctx context.Context, in *pb.UserId) (*pb.Profile, error) {
	profile := Profile{}
	if err := h.db.First(&profile, in.GetUserId()).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, fmt.Errorf("user does not exist")
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
	auth := Auth{}
	if err := h.db.Where("username = ?", in.GetUsername()).First(&auth).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(in.GetPassword())); err != nil {
		return nil, err
	}

	return &pb.UserId{UserId: proto.Uint64(auth.UserID)}, nil
}

func (h *Handlers) CreateAccount(ctx context.Context, in *pb.LoginCredentials) (*pb.UserId, error) {
	match, err := validUsername(in.GetUsername())
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, fmt.Errorf("invalid username")
	}

	match, err = validPassword(in.GetPassword())
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, fmt.Errorf("invalid password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, err
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
		return nil, err
	}

	profile := Profile{
		UserID:   auth.UserID,
		Username: auth.Username,
		Bio:      initialBio,
		AvatarID: initialAvatarID,
	}
	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &pb.UserId{UserId: proto.Uint64(auth.UserID)}, nil
}

func (h *Handlers) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	updates := map[string]interface{}{}
	updates["Bio"] = in.GetBio()
	updates["AvatarID"] = in.GetAvatarId()

	profile := Profile{UserID: in.GetUserId()}
	if err := h.db.Model(&profile).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &pb.UpdateProfileResponse{}, nil
}

func (h *Handlers) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, err
	}

	auth := Auth{}
	if err := tx.First(&auth, in.GetUserId()).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			return nil, fmt.Errorf("user does not exist")
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(in.GetOldPass())); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("invalid old password")
	}

	match, err := validPassword(in.GetNewPass())
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if !match {
		tx.Rollback()
		return nil, fmt.Errorf("invalid new password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.GetNewPass()), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	updates := map[string]interface{}{
		"password":   string(hash),
		"updated_at": time.Now(),
	}

	if err := tx.Model(&auth).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &pb.ChangePasswordResponse{}, nil
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
