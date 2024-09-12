package rest

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func (api *API) createToken(id int, email, authFor string) (string, time.Time, error) {
	exp_time, err := time.ParseDuration(api.Config.JwtExpires)
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt := time.Now().Add(exp_time)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
		"auth":  authFor,
		"exp":   expiresAt.Unix(), // Token expires in 24 hours
		"iat":   time.Now().Unix(),
		"type":  "access",
	})

	// SignedString expects a key for signing
	// You should replace "your-secret-key" with a actual secret key stored securely
	tokenString, err := token.SignedString([]byte(api.Config.JwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (api *API) createRefreshToken(id int, email, authFor string) (string, time.Time, error) {
	exp_time, err := time.ParseDuration(api.Config.RefreshExpiry)
	if err != nil {
		log.Println("error parsing refresh token expiry time", err, api.Config.RefreshExpiry)
		return "", time.Time{}, err
	}
	expiresAt := time.Now().Add(exp_time)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
		"auth":  authFor,
		"exp":   expiresAt.Unix(),  // Token expires in 7 days
		"iat":   time.Now().Unix(), // Issued at time
		"type":  "refresh",
	})

	// SignedString expects a key for signing
	// You should replace "your-secret-key" with a actual secret key stored securely
	tokenString, err := token.SignedString([]byte(api.Config.RefreshSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
