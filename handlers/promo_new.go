package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/99ridho/come-backend/errors"
	"github.com/99ridho/come-backend/models"
)

type newPromoRequest struct {
	Name          string  `json:"name"`
	Address       string  `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Description   string  `json:"description"`
	StartDate     string  `json:"start_date"`     // yyyy-MM-dd HH:mm
	EndDate       string  `json:"end_date"`       // yyyy-MM-dd HH:mm
	AllowedGender string  `json:"allowed_gender"` // male, female or both
	MaxSlot       int     `json:"max_slot"`
}

func NewPromo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req newPromoRequest

	if err := decoder.Decode(&req); err != nil {
		errors.NewErrorWithStatusCode(http.StatusBadRequest).WriteTo(w)
		return
	}

	userId := r.Context().Value("user_id").(int)
	if _, err := models.NewPromo(userId, req.Name, req.Address, req.Latitude, req.Longitude, req.Description, req.StartDate, req.EndDate, req.AllowedGender, req.MaxSlot); err != nil {
		errors.NewError("can't create new promo", http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "add promo success",
	})
}
