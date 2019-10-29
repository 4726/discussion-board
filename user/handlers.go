package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
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

type ErrorResponse struct {
	Error string
}

func GetProfile(db *gorm.DB, ctx *gin.Context) {
	useridS := ctx.Param("userid")

	userid, err := strconv.Atoi(useridS)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	profile := Profile{}
	if err := db.First(&profile, userid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, gin.H{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func ValidLogin(db *gorm.DB, ctx *gin.Context) {
	form := LoginForm{}
	err := ctx.BindJSON(&form)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	auth := Auth{}
	if err := db.Where("username = ?", form.Username).First(&auth).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(form.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{"invalid login"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"userid": auth.UserID})
}

func CreateAccount(db *gorm.DB, ctx *gin.Context) {
	form := CreateAccountForm{}
	err := ctx.BindJSON(&form)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	created := time.Now()
	auth := Auth{
		Username:  form.Username,
		Password:  string(hash),
		CreatedAt: created,
		UpdatedAt: created,
	}
	if err := db.Save(&auth).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"userid": auth.UserID})
}

func UpdateProfile(db *gorm.DB, ctx *gin.Context) {
	form := UpdateProfileForm{}
	err := ctx.BindJSON(&form)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if form.Bio != "" {
		updates["Bio"] = form.Bio
	}
	if form.AvatarID != "" {
		updates["AvatarID"] = form.AvatarID
	}

	profile := Profile{UserID: form.UserID}
	if err := db.Model(&profile).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func ChangePassword(db *gorm.DB, ctx *gin.Context) {
	form := ChangePasswordForm{}
	err := ctx.BindJSON(&form)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	auth := Auth{}
	if err := tx.First(&auth, form.UserID).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(form.OldPass)); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{"invalid old password"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(form.NewPass), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	if err := tx.Model(&auth).Update("password", string(hash)).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
