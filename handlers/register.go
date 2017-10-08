package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/99ridho/come-backend/models"
)

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Gender   string `json:"gender"`
	FcmToken string `json:"fcm_token"`
}

type registerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reqData registerRequest
	if err := decoder.Decode(&reqData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := models.NewUser(reqData.Username, reqData.Email, reqData.Password, reqData.FullName, reqData.Gender, reqData.FcmToken); err != nil {
		json.NewEncoder(w).Encode(registerResponse{
			Message: err.Error(),
			Status:  "failed",
		})
		return
	}

	json.NewEncoder(w).Encode(registerResponse{
		Message: "user registered",
		Status:  "success",
	})
}
