package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/99ridho/come-backend/errors"
	"github.com/99ridho/come-backend/models"
)

func ChangeProfilePhoto(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Photo string `json:"profile_photo"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		errors.NewErrorWithStatusCode(http.StatusBadRequest).WriteTo(w)
		return
	}

	userId := r.Context().Value("user_id").(int)

	var user models.User
	query := "SELECT * FROM users WHERE id = ?"
	if err := models.Dbm.SelectOne(&user, query, userId); err != nil {
		errors.NewErrorWithStatusCode(http.StatusInternalServerError).WriteTo(w)
		return
	}

	user.Photo = req.Photo
	if _, err := models.Dbm.Update(&user); err != nil {
		errors.NewError("can't update photo", http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "change photo succcessfully.",
	})
}
