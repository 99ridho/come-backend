package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"goji.io/pat"

	"github.com/99ridho/come-backend/errors"
	"github.com/99ridho/come-backend/models"
)

func DeleteMyPromoById(w http.ResponseWriter, r *http.Request) {
	promoId := pat.Param(r, "id")
	var promo models.Promo

	userId := r.Context().Value("user_id").(int)
	if err := models.Dbm.SelectOne(&promo, "select * from promos where id=? and user_id=?", promoId, userId); err != nil {
		errors.NewError("promo not found", http.StatusBadRequest).WriteTo(w)
		return
	}

	if err := promo.Delete(); err != nil {
		errors.NewError("error delete", http.StatusInternalServerError).WriteTo(w)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("delete promo with id %s successful", promoId),
	})
}
