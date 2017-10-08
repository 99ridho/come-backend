package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/99ridho/come-backend/env"
	"github.com/99ridho/come-backend/models"
	"github.com/dgrijalva/jwt-go"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Token   string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req loginRequest

	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.User
	query := "select * from users where email=?"

	if err := models.Dbm.SelectOne(&user, query, req.Email); err != nil {
		json.NewEncoder(w).Encode(loginResponse{
			Message: "user not found",
			Status:  "failed",
		})
		return
	}

	if err := user.VerifyPassword(req.Password); err != nil {
		json.NewEncoder(w).Encode(loginResponse{
			Message: "password incorrect",
			Status:  "failed",
		})
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
		Status:  "success",
		Token:   tokenString,
	})
}
