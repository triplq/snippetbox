package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              int
	Name            string
	Email           string
	Hashed_password []byte
	Created         time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	stmt := `INSERT INTO users(name, email, hash_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(stmt, name, email, hash)
	if err != nil {
		var mySqlError *mysql.MySQLError
		if errors.As(err, &mySqlError) {
			if mySqlError.Number == 1062 && strings.Contains(mySqlError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return true, nil
}
