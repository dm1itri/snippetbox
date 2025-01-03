package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID             int
	name           string
	email          string
	hashedPassword []byte
	created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (model *UserModel) Insert(name string, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (name, email, hashed_password, created) 
			  VALUES (?, ?, ?, UTC_TIMESTAMP())`
	_, err = model.DB.Exec(query, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (model *UserModel) Authenticate(email string, password string) (int, error) {
	var id int
	var hashedPassword []byte
	query := `SELECT id, hashed_password FROM users WHERE email = ?`
	err := model.DB.QueryRow(query, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (model *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
