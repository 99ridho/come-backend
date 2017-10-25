package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/99ridho/come-backend/errors"
	"github.com/99ridho/come-backend/models"
)

type registerRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	Gender      string `json:"gender"`
	FcmToken    string `json:"fcm_token"`
}

type registerResponse struct {
	Message string `json:"message"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reqData registerRequest
	if err := decoder.Decode(&reqData); err != nil {
		errors.NewErrorWithStatusCode(http.StatusBadRequest).WriteTo(w)
		return
	}

	if _, err := models.NewUser(reqData.Username, reqData.Email, reqData.PhoneNumber, reqData.Password, reqData.FullName, reqData.Gender, reqData.FcmToken, "user"); err != nil {
		errors.NewError("user already registered", http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(registerResponse{
		Message: "user registered",
	})
}
