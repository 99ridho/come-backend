package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"goji.io/pat"

	"github.com/99ridho/come-backend/errors"
	"github.com/99ridho/come-backend/models"
)

func ChangeMyPromoById(w http.ResponseWriter, r *http.Request) {
	promoId := pat.Param(r, "id")
	var req struct {
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
	var promo models.Promo

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		errors.NewErrorWithStatusCode(http.StatusBadRequest).WriteTo(w)
		return
	}

	userId := r.Context().Value("user_id").(int)
	if err := models.Dbm.SelectOne(&promo, "select * from promos where id=? and user_id=?", promoId, userId); err != nil {
		errors.NewError("promo not found", http.StatusBadRequest).WriteTo(w)
		return
	}

	if req.Address != "" {
		promo.Address = req.Address
	}
	if req.AllowedGender != "" {
		promo.AllowedGender = req.AllowedGender
	}
	if req.Description != "" {
		promo.Description = req.Description
	}
	if req.EndDate != "" {
		promo.EndDate = req.EndDate
	}
	if req.Latitude != 0 {
		promo.Latitude = req.Latitude
	}
	if req.Longitude != 0 {
		promo.Longitude = req.Longitude
	}
	if req.Name != "" {
		promo.Name = req.Name
	}
	if req.StartDate != "" {
		promo.StartDate = req.StartDate
	}
	promo.MaxSlot = req.MaxSlot
	if err := promo.Update(); err != nil {
		errors.NewError("error update", http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("update promo with id %s successful", promoId),
	})
}
