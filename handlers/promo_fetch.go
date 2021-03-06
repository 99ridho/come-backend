package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	p.*, u.full_name as owner_name, u.gender as owner_gender, u.role as owner_role,
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
	AND u.role = "user"
	AND (SELECT COUNT(*) FROM promo_attendees pa WHERE p.id = pa.promo_id) < p.max_slot
	HAVING distance < 5
	ORDER BY distance ASC
	`

	var promos []struct {
		models.Promo
		Distance    float64 `db:"distance" json:"distance"`
		OwnerName   string  `db:"owner_name" json:"owner_name"`
		OwnerGender string  `db:"owner_gender" json:"owner_gender"`
		OwnerRole   string  `db:"owner_role" json:"owner_role"`
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
		OwnerRole   string `db:"owner_role" json:"owner_role"`
	}

	query := "select p.*, u.full_name as owner_name, u.gender as owner_gender, u.role as owner_role from promos p, users u where p.id=? and u.id = p.user_id"
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
		OwnerRole   string `db:"owner_role" json:"owner_role"`
	}

	query := "select p.*, u.full_name as owner_name, u.gender as owner_gender, u.role as owner_role from promos p, users u where p.name LIKE ? AND u.id = p.user_id"
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

	query := `
	SELECT p.*, u.full_name as owner_name, u.gender as owner_gender, u.role as owner_role, u.phone_number as owner_phone_number
	FROM promos p
	INNER JOIN promo_attendees pa ON p.id = pa.promo_id
	INNER JOIN users u ON u.id = pa.user_id
	WHERE pa.user_id = ?
	`

	var promos []struct {
		models.Promo
		OwnerName        string `db:"owner_name" json:"owner_name"`
		OwnerGender      string `db:"owner_gender" json:"owner_gender"`
		OwnerPhoneNumber string `db:"owner_phone_number" json:"owner_phone_number"`
		OwnerRole        string `db:"owner_role" json:"owner_role"`
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
	promoId := pat.Param(r, "id")

	query := `
	SELECT u.id, u.username, u.email, u.phone_number, u.full_name, u.gender, u.fcm_token, u.role 
	FROM users u 
	INNER JOIN promo_attendees pa ON u.id = pa.user_id 
	INNER JOIN promos p ON pa.promo_id = p.id
	WHERE p.user_id = ? AND pa.promo_id = ?
	`

	var users []models.User
	if _, err := models.Dbm.Select(&users, query, myId, promoId); err != nil {
		errors.NewErrorWithStatusCode(http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": users,
	})
}

func FetchAdminPromotedPromo(w http.ResponseWriter, r *http.Request) {
	query := `
	SELECT p.*, u.full_name as owner_name, u.gender as owner_gender, u.role as owner_role 
	FROM promos p
	INNER JOIN users u ON u.id = p.user_id
	WHERE u.role = ?
	`

	var promos []struct {
		models.Promo
		OwnerName   string `db:"owner_name" json:"owner_name"`
		OwnerGender string `db:"owner_gender" json:"owner_gender"`
		OwnerRole   string `db:"owner_role" json:"owner_role"`
	}
	if _, err := models.Dbm.Select(&promos, query, "admin"); err != nil {
		errors.NewErrorWithStatusCode(http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": promos,
	})
}
