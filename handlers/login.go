package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/99ridho/come-backend/env"
	"github.com/99ridho/come-backend/errors"
	"github.com/99ridho/come-backend/models"
	"github.com/dgrijalva/jwt-go"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FcmToken string `json:"fcm_token"`
}

type loginResponse struct {
	Message      string      `json:"message"`
	Token        string      `json:"token"`
	LoggedInUser models.User `json:"logged_in_user"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req loginRequest

	if err := decoder.Decode(&req); err != nil {
		errors.NewErrorWithStatusCode(http.StatusBadRequest).WriteTo(w)
		return
	}

	var user models.User
	query := "select * from users where username=?"

	if err := models.Dbm.SelectOne(&user, query, req.Username); err != nil {
		errors.NewError("user not found", http.StatusUnauthorized).WriteTo(w)
		return
	}

	if err := user.VerifyPassword(req.Password); err != nil {
		errors.NewError("password incorrect", http.StatusUnauthorized).WriteTo(w)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := env.Getenv("SECRET_KEY", "secret")
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		errors.NewError("can't sign your token", http.StatusInternalServerError).WriteTo(w)
		return
	}

	user.FcmToken = req.FcmToken
	if _, err := models.Dbm.Update(&user); err != nil {
		errors.NewError("error logged in", http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(loginResponse{
		Message:      "logged in",
		Token:        tokenString,
		LoggedInUser: user,
	})
}
