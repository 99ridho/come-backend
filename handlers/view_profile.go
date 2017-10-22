package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"goji.io/pat"

	"github.com/99ridho/come-backend/errors"
	"github.com/99ridho/come-backend/models"
)

func FetchMyProfile(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("user_id").(int)

	user, err := fetchUser(strconv.Itoa(myId))

	if err != nil {
		errors.NewError("can't fetch profile", http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]models.User{
		"data": user,
	})
}

func FetchUserProfile(w http.ResponseWriter, r *http.Request) {
	userId := pat.Param(r, "id")

	user, err := fetchUser(userId)

	if err != nil {
		errors.NewError("can't fetch profile", http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]models.User{
		"data": user,
	})
}

func FetchUserProfileByUsername(w http.ResponseWriter, r *http.Request) {
	username := pat.Param(r, "username")

	var user models.User

	query := `
	SELECT id, username, email, full_name, gender, fcm_token
	FROM users
	WHERE username=?
	`
	if err := models.Dbm.SelectOne(&user, query, username); err != nil {
		errors.NewError("can't fetch profile", http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]models.User{
		"data": user,
	})
}

func fetchUser(id string) (models.User, error) {
	var user models.User

	query := `
	SELECT id, username, email, full_name, gender, fcm_token
	FROM users
	WHERE id=?
	`
	if err := models.Dbm.SelectOne(&user, query, id); err != nil {
		return models.User{}, err
	}

	return user, nil
}
