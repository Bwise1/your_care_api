package rest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

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

	token, token_expires, err := api.createToken(user.ID)
	if err != nil {
		return model.LoginResponse{}, values.Error, fmt.Sprintf("%s [CrTk]", values.SystemErr), err
	}

	refresh, refresh_expires, err := api.createRefreshToken(user.ID)
	if err != nil {
		log.Println(err)
		return model.LoginResponse{}, values.Error, fmt.Sprintf("%s [CrRF]", values.SystemErr), err
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

func (api *API) googleLogin() {

}

func (api *API) verifyEmail() {

}
