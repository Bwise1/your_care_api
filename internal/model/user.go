package model

import (
	"time"
)

// User represents the users table in the database
// type User struct {
// 	ID                            int             `json:"id"`
// 	FirstName                     string          `json:"firstName"`
// 	LastName                      string          `json:"lastName"`
// 	Email                         string          `json:"email"`
// 	DateOfBirth                   time.Time       `json:"dateOfBirth"`
// 	Sex                           string          `json:"sex"`
// 	Height                        sql.NullFloat64 `json:"height"` // Use sql.NullFloat64 for nullable fields
// 	Password                      string          `json:"-"`      // Exclude password from JSON output
// 	RoleID                        int             `json:"roleId"`
// 	IsActive                      bool            `json:"isActive"`
// 	LastLogin                     sql.NullTime    `json:"lastLogin"` // Use sql.NullTime for nullable datetime fields
// 	RefreshToken                  sql.NullString  `json:"-"`         // Exclude from JSON output
// 	TokenExpiration               sql.NullTime    `json:"tokenExpiration"`
// 	IsEmailVerified               bool            `json:"isEmailVerified"`
// 	EmailVerificationToken        sql.NullString  `json:"-"` // Exclude from JSON output
// 	EmailVerificationTokenExpires sql.NullTime    `json:"emailVerificationTokenExpires"`
// 	CreatedAt                     time.Time       `json:"createdAt"`
// 	UpdatedAt                     time.Time       `json:"updatedAt"`
// }

type User struct {
	ID                            int        `json:"id"`
	FirstName                     string     `json:"firstName"`
	LastName                      string     `json:"lastName"`
	Email                         string     `json:"email"`
	DateOfBirth                   time.Time  `json:"dateOfBirth"`
	Sex                           string     `json:"sex"`
	Height                        *float64   `json:"height,omitempty"`
	Password                      string     `json:"-"`
	RoleID                        int        `json:"roleId"`
	Role                          string     `json:"role"`
	IsActive                      bool       `json:"isActive"`
	LastLogin                     *time.Time `json:"lastLogin,omitempty"`
	RefreshToken                  *string    `json:"-"`
	TokenExpiration               *time.Time `json:"tokenExpiration,omitempty"`
	IsEmailVerified               bool       `json:"isEmailVerified"`
	EmailVerificationToken        *string    `json:"-"`
	EmailVerificationTokenExpires *time.Time `json:"emailVerificationTokenExpires,omitempty"`
	CreatedAt                     time.Time  `json:"createdAt"`
	UpdatedAt                     time.Time  `json:"updatedAt"`
}
type UserRequest struct {
	FirstName                    string    `json:"first_name"`
	LastName                     string    `json:"last_name"`
	Email                        string    `json:"email"`
	Password                     string    `json:"password"`
	DateOfBirth                  string    `json:"date_of_birth"`
	Sex                          string    `json:"sex"`
	EmailVerificationCode        string    `json:"-"`
	EmailVerificationCodeExpires time.Time `json:"-"`
}

type UserLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenInfo struct {
	AccessToken        string    `json:"access_token"`
	AccessTokenExpiry  time.Time `json:"access_token_expiry"`
	RefreshToken       string    `json:"refresh_token"`
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
}

type LoginResponse struct {
	User  User      `json:"user"`
	Token TokenInfo `json:"token"`
}

type EmailVerificationReq struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type LogotRequest struct {
	User  User      `json:"user"`
	Token TokenInfo `json:"token"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type ResendVerificationReq struct {
	Email string `json:"email" validate:"required,email"`
}

// model/user.go

type ChangePasswordReq struct {
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type EmailReq struct {
	Email string `json:"email" validate:"required,email"`
}
