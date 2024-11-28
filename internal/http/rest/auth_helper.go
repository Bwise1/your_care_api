package rest

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bwise1/your_care_api/internal/model"
	"github.com/bwise1/your_care_api/util"
	"github.com/bwise1/your_care_api/util/values"
	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	UserID int    `json:"sub"`
	Type   string `json:"typ"`
	Exp    int64  `json:"exp"`
}

// Simplified token creation
func (api *API) createToken(id int) (string, time.Time, error) {
	exp_time, err := time.ParseDuration(api.Config.JwtExpires)
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt := time.Now().Add(exp_time)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id, // subject (user ID)
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
		"typ": "access",
	})

	tokenString, err := token.SignedString([]byte(api.Config.JwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (api *API) createRefreshToken(id int) (string, time.Time, error) {
	exp_time, err := time.ParseDuration(api.Config.RefreshExpiry)
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt := time.Now().Add(exp_time)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id, // subject (user ID)
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
		"typ": "refresh",
	})

	tokenString, err := token.SignedString([]byte(api.Config.RefreshSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (api *API) verifyToken(tokenString string, isRefresh bool) (*TokenClaims, error) {

	// Determine the correct secret key based on token type
	secret := api.Config.JwtSecret
	if isRefresh {
		secret = api.Config.RefreshSecret
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	// Specifically handle token expiration
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, fmt.Errorf("token expired")
		}
	}

	// Check for errors or invalid token
	if err != nil || !token.Valid {
		log.Println("error verifying token", err)
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	// Check the token type (use "typ" instead of "type")
	tokenType, _ := claims["typ"].(string) // Fixed to check "typ" claim
	if (isRefresh && tokenType != "refresh") || (!isRefresh && tokenType != "access") {
		return nil, fmt.Errorf("invalid token type")
	}

	// Extract user ID
	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	userID := int(userIDFloat)

	// Verify the refresh token from the database if it is a refresh token
	if isRefresh {
		log.Println("verifying refresh token in database")
		valid, err := api.verifyRefreshTokenInDB(context.TODO(), userID, tokenString)
		if err != nil {
			return nil, fmt.Errorf("error verifying refresh token from database: %v", err)
		}
		if !valid {
			return nil, fmt.Errorf("invalid refresh token")
		}
	}

	// Log extracted user ID and token type
	log.Println("user id", userID)
	log.Println("token type", tokenType)

	// Return the extracted claims
	return &TokenClaims{
		UserID: userID,
		Type:   tokenType,
		Exp:    int64(claims["exp"].(float64)),
	}, nil
}

func (api *API) LogUserOut(userID int) (bool, error) {
	err := api.invalidateRefreshToken(context.TODO(), userID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// auth_helper.go

func (api *API) ResendVerificationEmail(req model.ResendVerificationReq) (string, string, error) {
	// Validate email format
	if err := util.ValidEmail(req.Email); err != nil {
		return values.BadRequestBody, "invalid email format", err
	}

	// Get user by email
	user, err := api.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return values.NotFound, "user not found", fmt.Errorf("user not found")
		}
		return values.Error, "error getting user", err
	}

	// Check if email is already verified
	if user.IsEmailVerified {
		return values.BadRequestBody, "email already verified", fmt.Errorf("email already verified")
	}

	// Generate new verification code
	verificationCode := util.RandomString(6, values.Numbers)
	expiryTime := time.Now().Add(time.Minute * 10)

	// Update user's verification code and expiry
	err = api.updateVerificationCode(context.Background(), user.ID, verificationCode, expiryTime)
	if err != nil {
		return values.Error, "error updating verification code", err
	}

	//TODO: Send verification email
	// err = api.Deps.EmailService.SendVerificationEmail(user.Email, user.FirstName, verificationCode)
	// if err != nil {
	// 	return values.Error, "error sending verification email", err
	// }

	return values.Success, "verification email sent successfully", nil
}

// func (api *API) verifyEmail(emailReq model.EmailVerificationReq) (bool, error) {
// 	emailReq.Email = strings.Trim(emailReq.Email, " ")

// 	err := util.ValidEmail(emailReq.Email)
// 	if err != nil {
// 		return false, err
// 	}

// 	//confirm if the users code rhymes with the email veridicatuin code also check if it has expired

// 	user, err := api.GetUserByEmail(context.TODO(), emailReq.Email)
// 	if err != nil {
// 		if err.Error() == values.NotFound {
// 			return false, fmt.Errorf("user does not exist")
// 		}
// 		return false, err
// 	}
// 	if user.IsEmailVerified {
// 		return false, fmt.Errorf("email already verified")
// 	}

// 	if !user {
// 		return false, fmt.Errorf("no email verification code")
// 	}

// }
