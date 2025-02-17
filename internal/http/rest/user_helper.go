package rest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwise1/your_care_api/internal/model"
	"github.com/bwise1/your_care_api/util"
	"github.com/bwise1/your_care_api/util/values"
	"golang.org/x/crypto/bcrypt"
)

func (api *API) RegisterUser(req model.UserRequest) (model.User, string, string, error) {

	var err error
	var ctx = context.TODO()

	req.Email = strings.Trim(req.Email, " ")
	req.FirstName = strings.Trim(req.FirstName, " ")
	req.LastName = strings.Trim(req.LastName, " ")

	err = util.ValidEmail(req.Email)
	if err != nil {
		return model.User{}, values.NotAllowed, "Invalid email address provided", err
	}

	exists, err := api.EmailExists(ctx, req.Email)
	if err != nil {
		return model.User{}, values.Error, fmt.Sprintf("%s [EmCh]", values.SystemErr), err
	}
	if exists {
		return model.User{}, values.Conflict, "User already exists. Please login", errors.New(values.Conflict)
	}

	passHash, err := util.HashPassword([]byte(req.Password))
	if err != nil {
		return model.User{}, values.Error, fmt.Sprintf("%s [HsPw]", values.SystemErr), err
	}

	req.Password = passHash

	verificationCode := util.RandomString(6, values.Numbers)

	verificationCodeExpires := time.Now().Add(time.Minute * 10)

	req.EmailVerificationCode = verificationCode
	req.EmailVerificationCodeExpires = verificationCodeExpires

	data := struct {
		Name            string
		VerificationURL string
	}{
		Name:            req.FirstName,
		VerificationURL: "http://localhost:3000?token=" + verificationCode,
	}
	patterns := []string{"verifyEmail.tmpl"}
	err = api.Deps.Mailer.Send(req.Email, data, patterns...)
	if err != nil {
		log.Println(err)
	}

	err = api.CreateUserRepo(context.TODO(), req)
	if err != nil {
		return model.User{}, values.Error, fmt.Sprintf("%s [CrUs]", values.SystemErr), err
	}

	user := model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}
	return user, values.Created, "Registration completed successfully", nil
}

func (api *API) LoginUser(req model.UserLoginReq) (model.LoginResponse, string, string, error) {

	var err error
	var ctx = context.TODO()

	req.Email = strings.Trim(req.Email, " ")

	err = util.ValidEmail(req.Email)
	if err != nil {
		return model.LoginResponse{}, values.NotAllowed, "Invalid email address provided", err
	}

	exists, err := api.EmailExists(ctx, req.Email)
	if err != nil {
		return model.LoginResponse{}, values.Error, fmt.Sprintf("%s [EmCh]", values.SystemErr), err
	}
	if !exists {
		return model.LoginResponse{}, values.NotFound, "User does not exist. Please register", errors.New(values.NotFound)
	}

	user, err := api.GetUserByEmail(ctx, req.Email)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return model.LoginResponse{}, values.NotFound, "User does not exist", err
		}
		return model.LoginResponse{}, values.Error, fmt.Sprintf("%s [GtUs]", values.SystemErr), err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return model.LoginResponse{}, values.NotAuthorised, "Invalid password provided", err
	}

	token, token_expires, err := api.createToken(user.ID, user.Role)
	if err != nil {
		return model.LoginResponse{}, values.Error, fmt.Sprintf("%s [CrTk]", values.SystemErr), err
	}

	refresh, refresh_expires, err := api.createRefreshToken(user.ID)
	if err != nil {
		log.Println(err)
		return model.LoginResponse{}, values.Error, fmt.Sprintf("%s [CrRF]", values.SystemErr), err
	}

	err = api.StoreRefreshToken(ctx, user.ID, refresh, refresh_expires)
	log.Println("user id", user.ID)
	if err != nil {
		return model.LoginResponse{}, values.Error, fmt.Sprintf("%s [StRF]", values.SystemErr), err
	}

	loggedUser := model.LoginResponse{
		User: user,
		Token: model.TokenInfo{
			AccessToken:        token,
			AccessTokenExpiry:  token_expires,
			RefreshToken:       refresh,
			RefreshTokenExpiry: refresh_expires,
		},
	}
	return loggedUser, values.Success, "User authenticated successfully", nil
}

func (api *API) RefreshToken(req model.RefreshTokenReq) (model.TokenInfo, string, string, error) {

	var err error
	var ctx = context.TODO()

	tokenClaims, err := api.verifyToken(req.RefreshToken, true)

	if err != nil {
		return model.TokenInfo{}, values.NotAuthorised, "Invalid refresh token", err
	}

	newAccessToken, accessTokenExpiry, err := api.createToken(tokenClaims.UserID, tokenClaims.Role)
	if err != nil {
		return model.TokenInfo{}, values.Error, "Failed to create new access token", err
	}

	newRefreshToken, refreshTokenExpiry, err := api.createRefreshToken(tokenClaims.UserID)
	if err != nil {
		return model.TokenInfo{}, values.Error, "Failed to create new refresh token", err
	}

	err = api.StoreRefreshToken(ctx, tokenClaims.UserID, newRefreshToken, refreshTokenExpiry)
	if err != nil {
		return model.TokenInfo{}, values.Error, "Failed to store new refresh token", err
	}

	return model.TokenInfo{
		AccessToken:        newAccessToken,
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshToken:       newRefreshToken,
		RefreshTokenExpiry: refreshTokenExpiry,
	}, values.Success, "Token refreshed successfully", nil
}

func (api *API) ChangeUserPassword(userID int, req model.ChangePasswordReq) (string, string, error) {
	// Get current user
	user, err := api.GetUserByID(context.Background(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return values.NotFound, "user not found", err
		}
		return values.Error, "error getting user", err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return values.BadRequestBody, "invalid old password", err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return values.Error, "error hashing password", err
	}

	// Update password in database
	err = api.UpdateUserPassword(context.Background(), userID, string(hashedPassword))
	if err != nil {
		return values.Error, "error updating password", err
	}

	return values.Success, "password updated successfully", nil
}

func (api *API) CheckEmail(req model.EmailReq) (string, string, error) {
	// Validate email format
	if err := util.ValidEmail(req.Email); err != nil {
		return values.BadRequestBody, "invalid email format", err
	}

	// Check if email exists
	exists, err := api.EmailExists(context.Background(), req.Email)
	if err != nil {
		return values.Error, "error checking email", err
	}
	if !exists {
		return values.NotFound, "email not found", nil
	}

	return values.Success, "email found", nil
}

func (api *API) googleLogin() {

}

func (api *API) RegisterAdminUser(req model.UserRequest) (model.User, string, string, error) {
	var err error
	var ctx = context.TODO()

	req.Email = strings.Trim(req.Email, " ")
	req.FirstName = strings.Trim(req.FirstName, " ")
	req.LastName = strings.Trim(req.LastName, " ")

	err = util.ValidEmail(req.Email)
	if err != nil {
		return model.User{}, values.NotAllowed, "Invalid email address provided", err
	}

	exists, err := api.EmailExists(ctx, req.Email)
	if err != nil {
		return model.User{}, values.Error, fmt.Sprintf("%s [EmCh]", values.SystemErr), err
	}
	if exists {
		return model.User{}, values.Conflict, "User already exists. Please login", errors.New(values.Conflict)
	}

	passHash, err := util.HashPassword([]byte(req.Password))
	if err != nil {
		return model.User{}, values.Error, fmt.Sprintf("%s [HsPw]", values.SystemErr), err
	}

	req.Password = passHash

	// verificationCode := util.RandomString(6, values.Numbers)
	// verificationCodeExpires := time.Now().Add(time.Minute * 10)

	// req.EmailVerificationCode = verificationCode
	// req.EmailVerificationCodeExpires = verificationCodeExpires

	// data := struct {
	//     Name            string
	//     VerificationURL string
	// }{
	//     Name:            req.FirstName,
	//     VerificationURL: "http://localhost:3000?token=" + verificationCode,
	// }
	// patterns := []string{"verifyEmail.tmpl"}
	// err = api.Deps.Mailer.Send(req.Email, data, patterns...)
	// if err != nil {
	//     log.Println(err)
	// }

	err = api.CreateAdminUserRepo(context.TODO(), req)
	if err != nil {
		return model.User{}, values.Error, fmt.Sprintf("%s [CrUs]", values.SystemErr), err
	}

	user := model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}
	return user, values.Created, "Admin registration completed successfully", nil
}
