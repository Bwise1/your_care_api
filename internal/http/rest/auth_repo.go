package rest

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bwise1/your_care_api/internal/model"
)

func (api *API) CreateUserRepo(ctx context.Context, req model.UserRequest) error {
	log.Println("creating user, ", req)

	stmt := `INSERT INTO users(
		firstName,
		lastName,
		email,
		password,
		dateOfBirth,
		sex,
		emailVerificationCode,
		emailVerificationCodeExpires
	)VALUES(?, ?, ?, ?, ?, ?,?,?)`

	_, err := api.Deps.DB.ExecContext(ctx, stmt, req.FirstName, req.LastName, req.Email, req.Password, req.DateOfBirth, req.Sex, req.EmailVerificationCode, req.EmailVerificationCodeExpires)
	if err != nil {
		log.Println("error creating user", err)
		return err
	}
	return nil
}

func (api *API) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`

	err := api.Deps.DB.QueryRowContext(ctx, stmt, email).Scan(&exists)
	if err != nil {
		log.Println("error checking email", err)
		return false, err
	}
	return exists, nil
}

func (api *API) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	stmt := `SELECT
		id,
		firstName,
		lastName,
		email,
		password,
		role_id,
		isActive,
		isEmailVerified
	FROM users
	WHERE email = ?`

	err := api.Deps.DB.QueryRowContext(ctx, stmt, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.IsActive,
		&user.IsEmailVerified,
	)
	if err != nil {
		log.Println("error getting user", err)
		return model.User{}, err
	}
	return user, nil
}

func (api *API) GetUserByID(ctx context.Context, id int) (model.User, error) {
	var user model.User
	stmt := `SELECT
		id,
		firstName,
		lastName,
		email,
		password,
		role_id,
		isActive,
		isEmailVerified
	FROM users
	WHERE id = ?`

	err := api.Deps.DB.QueryRowContext(ctx, stmt, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.IsActive,
		&user.IsEmailVerified,
	)
	if err != nil {
		log.Println("error getting user", err)
		return model.User{}, err
	}
	return user, nil
}

func (api *API) GetUserBySocialID(provider, socialID string) (*model.User, error) {
	// Implement database query to find user by social ID
	// This is a placeholder function
	return nil, nil
}

func (api *API) StoreRefreshToken(ctx context.Context, userID int, refreshToken string, expiry time.Time) error {
	stmt := `UPDATE users SET refreshToken = ? WHERE id = ?`
	_, err := api.Deps.DB.ExecContext(ctx, stmt, refreshToken, userID)
	if err != nil {
		log.Println("error storing refresh token", err)
		return err
	}
	return nil
}

func (api *API) invalidateRefreshToken(ctx context.Context, userID int) error {
	stmt := `UPDATE users SET refresh_token = NULL WHERE id = ?`
	_, err := api.Deps.DB.ExecContext(ctx, stmt, userID)
	if err != nil {
		log.Println("error invalidating refresh token", err)
		return err
	}
	return nil
}

func (api *API) verifyRefreshTokenInDB(ctx context.Context, userID int, refreshToken string) (bool, error) {
	var token string
	stmt := `SELECT refreshToken FROM users WHERE id = ?`

	err := api.Deps.DB.QueryRowContext(ctx, stmt, userID).Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error verifying refresh token: %w", err)
	}

	// Use constant-time comparison
	return subtle.ConstantTimeCompare([]byte(token), []byte(refreshToken)) == 1, nil
}

func (api *API) updateVerificationCode(ctx context.Context, userID int, code string, expiry time.Time) error {
	stmt := `UPDATE users SET
        emailVerificationCode = ?,
        emailVerificationCodeExpires = ?
    WHERE id = ?`

	_, err := api.Deps.DB.ExecContext(ctx, stmt, code, expiry, userID)
	if err != nil {
		log.Println("error updating verification code:", err)
		return err
	}

	return nil

}
func (api *API) UpdateUserPassword(ctx context.Context, userID int, hashedPassword string) error {
	stmt := `UPDATE users SET password = ? WHERE id = ?`
	result, err := api.Deps.DB.ExecContext(ctx, stmt, hashedPassword, userID)
	if err != nil {
		log.Println("error updating password:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID: %d", userID)
	}

	return nil
}

func (api *API) CreateAdminUserRepo(ctx context.Context, req model.UserRequest) error {
	log.Println("creating admin user, ", req)

	stmt := `INSERT INTO users(
        firstName,
        lastName,
        email,
        password,
        sex,
        role_id
    )VALUES(?, ?, ?, ?, ?, ?, (SELECT id FROM roles WHERE name = 'admin'), ?, ?)`

	_, err := api.Deps.DB.ExecContext(ctx, stmt, req.FirstName, req.LastName, req.Email, req.Password, "Other")
	if err != nil {
		log.Println("error creating admin user", err)
		return err
	}
	return nil
}
