package handlers

import (
	"encoding/json"
	"net/http"

	"goji.io/pat"

	"github.com/99ridho/come-backend/errors"
	"github.com/99ridho/come-backend/models"
)

func DeletePromoAttendee(w http.ResponseWriter, r *http.Request) {
	attendeeId := pat.Param(r, "attendee_id")
	promoId := pat.Param(r, "promo_id")

	query := "DELETE FROM promo_attendees WHERE user_id = ? AND promo_id = ?"
	if _, err := models.Dbm.Exec(query, attendeeId, promoId); err != nil {
		errors.NewError("can't fetch attendee", http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "delete attendee successfully.",
	})
}
