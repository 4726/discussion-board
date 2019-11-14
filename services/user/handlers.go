package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type ErrorResponse struct {
	Error string
}

var (
	InvalidJSONBodyResponse = ErrorResponse{"invalid body"}
)

const (
	initialBio, initialAvatarID = "", ""
)

func GetProfile(db *gorm.DB, ctx *gin.Context) {
	useridS := ctx.Param("userid")

	userid, err := strconv.Atoi(useridS)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid userid param"})
		return
	}

	profile := Profile{}
	if err := db.First(&profile, userid).Error; err != nil {
		ctx.Set(logInfoKey, err)
		if gorm.IsRecordNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, struct{}{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func ValidLogin(db *gorm.DB, ctx *gin.Context) {
	form := struct {
		Username string `binding:"required"`
		Password string `binding:"required"`
	}{}
	err := ctx.BindJSON(&form)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	auth := Auth{}
	if err := db.Where("username = ?", form.Username).First(&auth).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			ctx.Set(logInfoKey, err)
			ctx.JSON(http.StatusUnauthorized, ErrorResponse{"invalid login"})
			return
		}
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(form.Password)); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{"invalid login"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"userid": auth.UserID})
}

func CreateAccount(db *gorm.DB, ctx *gin.Context) {
	form := struct {
		Username string `binding:"required"`
		Password string `binding:"required"`
	}{}
	err := ctx.BindJSON(&form)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	match, err := validUsername(form.Username)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	if !match {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid username"})
		return
	}

	match, err = validPassword(form.Password)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	if !match {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid password"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	created := time.Now()
	auth := Auth{
		Username:  form.Username,
		Password:  string(hash),
		CreatedAt: created,
		UpdatedAt: created,
	}
	if err := tx.Save(&auth).Error; err != nil {
		tx.Rollback()
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	profile := Profile{
		UserID:   auth.UserID,
		Username: auth.Username,
		Bio:      initialBio,
		AvatarID: initialAvatarID,
	}
	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"userid": auth.UserID})
}

func UpdateProfile(db *gorm.DB, ctx *gin.Context) {
	form := struct {
		UserID        uint `binding:"required"`
		Bio, AvatarID string
	}{}
	err := ctx.BindJSON(&form)
	if err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
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
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, struct{}{})
}

func ChangePassword(db *gorm.DB, ctx *gin.Context) {
	form := struct {
		UserID  int    `binding:"required"`
		OldPass string `binding:"required"`
		NewPass string `binding:"required"`
	}{}
	if err := ctx.BindJSON(&form); err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, InvalidJSONBodyResponse)
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	auth := Auth{}
	if err := tx.First(&auth, form.UserID).Error; err != nil {
		tx.Rollback()
		if gorm.IsRecordNotFoundError(err) {
			ctx.Set(logInfoKey, err)
			ctx.JSON(http.StatusNotFound, struct{}{})
			return
		}
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(form.OldPass)); err != nil {
		tx.Rollback()
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{"invalid old password"})
		return
	}

	match, err := validPassword(form.NewPass)
	if err != nil {
		tx.Rollback()
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	if !match {
		tx.Rollback()
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusBadRequest, ErrorResponse{"invalid new password"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(form.NewPass), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	updates := map[string]interface{}{
		"password":   string(hash),
		"updated_at": time.Now(),
	}

	if err := tx.Model(&auth).Updates(updates).Error; err != nil {
		tx.Rollback()
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.Set(logInfoKey, err)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{"server error"})
		return
	}

	ctx.JSON(http.StatusOK, struct{}{})
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
