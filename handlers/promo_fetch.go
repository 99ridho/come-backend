package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"goji.io/pat"

	"github.com/99ridho/come-backend/errors"
	"github.com/99ridho/come-backend/models"
)

func FetchMyPromos(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(int)
	var promos []models.Promo

	if _, err := models.Dbm.Select(&promos, "select * from promos where user_id=?", userId); err != nil {
		errors.NewErrorWithStatusCode(http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string][]models.Promo{
		"data": promos,
	})
}

func FetchNearbyPromosFromMyLocation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		errors.NewErrorWithStatusCode(http.StatusBadRequest).WriteTo(w)
		return
	}

	// lat, lng, lat
	query := `
	SELECT 
	p.*, u.full_name as owner_name, u.gender as owner_gender,
	(
	   6371 *
	   acos(cos(radians(?)) * 
	   cos(radians(p.latitude)) * 
	   cos(radians(p.longitude) - 
	   radians(?)) + 
	   sin(radians(?)) * 
	   sin(radians(p.latitude)))
	) AS distance 
	FROM promos p, users u
	WHERE u.id = p.user_id
	HAVING distance < 5
	ORDER BY distance ASC
	`

	var promos []struct {
		models.Promo
		Distance    float64 `db:"distance" json:"distance"`
		OwnerName   string  `db:"owner_name" json:"owner_name"`
		OwnerGender string  `db:"owner_gender" json:"owner_gender"`
	}
	if _, err := models.Dbm.Select(&promos, query, req.Latitude, req.Longitude, req.Latitude); err != nil {
		errors.NewError("can't fetch nearby promo", http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": promos,
	})
}

func FetchPromoById(w http.ResponseWriter, r *http.Request) {
	promoId := pat.Param(r, "id")
	//var promo models.Promo

	var promo struct {
		models.Promo
		OwnerName   string `db:"owner_name" json:"owner_name"`
		OwnerGender string `db:"owner_gender" json:"owner_gender"`
	}

	query := "select p.*, u.full_name as owner_name, u.gender as owner_gender from promos p, users u where p.id=? and u.id = p.user_id"
	if err := models.Dbm.SelectOne(&promo, query, promoId); err != nil {
		errors.NewError("can't fetch promo", http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": promo,
	})
}

func FetchPromoByName(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Query string `json:"query"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		errors.NewErrorWithStatusCode(http.StatusBadRequest).WriteTo(w)
		return
	}

	var promos []struct {
		models.Promo
		OwnerName   string `db:"owner_name" json:"owner_name"`
		OwnerGender string `db:"owner_gender" json:"owner_gender"`
	}

	query := "select p.*, u.full_name as owner_name, u.gender as owner_gender from promos p, users u where p.name LIKE ? AND u.id = p.user_id"
	if _, err := models.Dbm.Select(&promos, query, fmt.Sprintf("%%%s%%", req.Query)); err != nil {
		errors.NewError(err.Error(), http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": promos,
	})
}

func FetchJoinedPromo(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("user_id").(int)

	query := "select p.*, u.full_name as owner_name, u.gender as owner_gender from promos p, users u, promo_attendees pa where pa.user_id=? and u.id = p.user_id and pa.promo_id = p.id"

	var promos []struct {
		models.Promo
		OwnerName   string `db:"owner_name" json:"owner_name"`
		OwnerGender string `db:"owner_gender" json:"owner_gender"`
	}

	if _, err := models.Dbm.Select(&promos, query, myId); err != nil {
		errors.NewError(err.Error(), http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": promos,
	})
}

func FetchMyPromoAttendeeByPromoId(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("user_id").(int)
	promoId, _ := strconv.Atoi(pat.Param(r, "id"))

	// query := "select * from users u, promo_attendees pa, promos p where u.id = pa.user_id and p.user_id = ? and pa.promo_id = p.id and p.id = ?"
	query := "select * from promo_attendees"
	var attendees []models.PromoAttendee
	users := make([]*models.User, 0)

	if _, err := models.Dbm.Select(&attendees, query); err != nil {
		errors.NewError(err.Error(), http.StatusInternalServerError).WriteTo(w)
		return
	}

	for _, a := range attendees {
		promo, _ := a.Promo()
		if promo.UserID == myId && promo.ID == promoId {
			user, _ := a.User()
			users = append(users, user)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": users,
	})
}
