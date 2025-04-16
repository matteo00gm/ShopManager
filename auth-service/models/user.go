package models

import (
	"database/sql"
	"errors"

	"shopsweb.com/auth-service/db"
	"shopsweb.com/auth-service/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := "INSERT INTO users(email, password) VALUES($1, $2) RETURNING id"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	// Use QueryRow to get the generated ID
	err = stmt.QueryRow(u.Email, hashedPassword).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = $1"
	row := db.DB.QueryRow(query, u.Email)
	var retrievedPassword string

	//mapping results to userId and retrievedPassword
	err := row.Scan(&u.ID, &retrievedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//no rows found, but i'm using this message to avoid users guessing which part of the login failed
			return errors.New("credentials invalid")
		}
		//other errors
		return err
	}

	// Check if the provided password matches the hashed password
	passIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)

	if !passIsValid {
		return errors.New("credentials invalid")
	}

	return nil
}
