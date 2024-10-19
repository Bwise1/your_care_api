package rest

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

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

// Simplified claims structure
type TokenClaims struct {
	UserID int    `json:"sub"`
	Type   string `json:"typ"`
	Exp    int64  `json:"exp"`
}

// Simplified token verification
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

	// Log the claims for debugging
	log.Println("claims", claims)

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
