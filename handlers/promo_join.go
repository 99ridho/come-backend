package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"goji.io/pat"

	"github.com/99ridho/come-backend/errors"
	"github.com/99ridho/come-backend/models"
)

func JoinPromo(w http.ResponseWriter, r *http.Request) {
	promoId, _ := strconv.ParseInt(pat.Param(r, "id"), 10, 0)

	// 1. select user
	userId := r.Context().Value("user_id").(int)
	var user models.User
	if err := models.Dbm.SelectOne(&user, "select * from users where id=?", userId); err != nil {
		errors.NewErrorWithStatusCode(http.StatusInternalServerError).WriteTo(w)
		return
	}

	// 2. count promo attendee
	attendeeCount, countError := models.Dbm.SelectInt("select count(*) from promo_attendees where promo_id=?", promoId)
	if countError != nil {
		errors.NewErrorWithStatusCode(http.StatusInternalServerError).WriteTo(w)
		return
	}
	var promo models.Promo
	if err := models.Dbm.SelectOne(&promo, "select * from promos where id=?", promoId); err != nil {
		errors.NewErrorWithStatusCode(http.StatusInternalServerError).WriteTo(w)
		return
	}

	// 3. check allowed gender
	if promo.AllowedGender != "both" {
		if promo.AllowedGender != user.Gender {
			errors.NewError("can't join promo, gender mismatch!", http.StatusInternalServerError).WriteTo(w)
			return
		}
	}

	// 4. check if promo is full slot
	if int(attendeeCount) >= promo.MaxSlot {
		errors.NewError("can't join promo, slot full!", http.StatusInternalServerError).WriteTo(w)
		return
	}

	// 5. check if inserting to DB encountered an error
	if _, err := models.NewPromoAttendee(userId, int(promoId)); err != nil {
		errors.NewError("can't join promo", http.StatusInternalServerError).WriteTo(w)
		return
	}

	// 6. TODO - notify promo owner..

	json.NewEncoder(w).Encode(map[string]string{
		"message": "joining promo success",
	})
}
