package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               int     `db:"id" json:"id"`
	Username         string  `db:"username" json:"username"`
	Email            string  `db:"email" json:"email"`
	PhoneNumber      string  `db:"phone_number" json:"phone_number"`
	Password         string  `db:"password" json:"-"`
	FullName         string  `db:"full_name" json:"full_name"`
	Gender           string  `db:"gender" json:"gender"`
	FcmToken         string  `db:"fcm_token" json:"fcm_token"`
	Role             string  `db:"role" json:"-"` // role: admin or user
	CurrentLatitude  float64 `db:"current_latitude" json:"current_latitude"`
	CurrentLongitude float64 `db:"current_longitude" json:"current_longitude"`
}

func NewUser(username string, email string, phone string, password string, fullName string, gender string, fcmToken string, role string, lat float64, lon float64) (*User, error) {
	user := &User{
		Username:         username,
		Email:            email,
		PhoneNumber:      phone,
		FullName:         fullName,
		Gender:           gender,
		FcmToken:         fcmToken,
		Role:             role,
		CurrentLatitude:  lat,
		CurrentLongitude: lon,
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
