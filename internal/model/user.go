package model

import (
	"database/sql"
	"time"
)

// User represents the users table in the database
type User struct {
	ID                            int             `json:"id"`
	FirstName                     string          `json:"firstName"`
	LastName                      string          `json:"lastName"`
	Email                         string          `json:"email"`
	DateOfBirth                   time.Time       `json:"dateOfBirth"`
	Sex                           string          `json:"sex"`
	Height                        sql.NullFloat64 `json:"height"` // Use sql.NullFloat64 for nullable fields
	Password                      string          `json:"-"`      // Exclude password from JSON output
	RoleID                        int             `json:"roleId"`
	IsActive                      bool            `json:"isActive"`
	LastLogin                     sql.NullTime    `json:"lastLogin"` // Use sql.NullTime for nullable datetime fields
	RefreshToken                  sql.NullString  `json:"-"`         // Exclude from JSON output
	TokenExpiration               sql.NullTime    `json:"tokenExpiration"`
	IsEmailVerified               bool            `json:"isEmailVerified"`
	EmailVerificationToken        sql.NullString  `json:"-"` // Exclude from JSON output
	EmailVerificationTokenExpires sql.NullTime    `json:"emailVerificationTokenExpires"`
	CreatedAt                     time.Time       `json:"createdAt"`
	UpdatedAt                     time.Time       `json:"updatedAt"`
}
