package userdao

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (user *UserDAO) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (user *UserDAO) HashPassword() error {
	if !user.Password.Valid {
		return errors.New("Invalid password")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password.String), 14)
	if err != nil {
		return err
	}
	user.Password.SetValid(string(bytes))
	return nil
}
