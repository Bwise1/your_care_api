package rest

import (
	"context"
	"log"

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
		sex
	)VALUES(?, ?, ?, ?, ?, ?)`

	_, err := api.Deps.DB.ExecContext(ctx, stmt, req.FirstName, req.LastName, req.Email, req.Password, req.DateOfBirth, req.Sex)
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
