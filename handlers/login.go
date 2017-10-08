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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req loginRequest

	if err := decoder.Decode(&req); err != nil {
		errors.NewErrorWithStatusCode(http.StatusBadRequest).WriteTo(w)
		return
	}

	var user models.User
	query := "select * from users where email=?"

	if err := models.Dbm.SelectOne(&user, query, req.Email); err != nil {
		errors.NewError("user not found", http.StatusInternalServerError).WriteTo(w)
		return
	}

	if err := user.VerifyPassword(req.Password); err != nil {
		errors.NewError("password incorrect", http.StatusInternalServerError).WriteTo(w)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := env.Getenv("SECRET_KEY", "secret")
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(loginResponse{
		Message: "logged in",
		Token:   tokenString,
	})
}
