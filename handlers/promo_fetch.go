package handlers

import (
	"encoding/json"
	"net/http"

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
	*, 
	(
	   6371 *
	   acos(cos(radians(?)) * 
	   cos(radians(latitude)) * 
	   cos(radians(longitude) - 
	   radians(?)) + 
	   sin(radians(?)) * 
	   sin(radians(latitude)))
	) AS distance 
	FROM promos
	HAVING distance < 5
	ORDER BY distance ASC
	`

	var promos []struct {
		models.Promo
		Distance float64 `json:"distance"`
	}
	if _, err := models.Dbm.Select(&promos, query, req.Latitude, req.Longitude, req.Latitude); err != nil {
		errors.NewError(err.Error(), http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": promos,
	})
}
