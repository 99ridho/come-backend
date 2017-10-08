package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
	FullName string `db:"full_name"`
	Gender   string `db:"gender"`
	FcmToken string `db:"fcm_token"`
}

func NewUser(id int, username string, email string, password string, fullName string, gender string, fcmToken string) (*User, error) {
	user := &User{
		ID:       id,
		Username: username,
		Email:    email,
		FullName: fullName,
		Gender:   gender,
		FcmToken: fcmToken,
	}
	user.HashPassword(password)
	err := Dbm.Insert(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) HashPassword(raw string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	u.Password = string(hashedPassword)
}

func (u *User) VerifyPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
