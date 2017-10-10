package models

import (
	"testing"

	"github.com/99ridho/come-backend/models"
)

func TestIfPromoRelationToUserIsCorrect(t *testing.T) {
	userId := 1 // example
	var promos []models.Promo
	if _, err := models.Dbm.Select(&promos, "select * from promos where user_id=?", userId); err != nil {
		t.Errorf("can't select data, expected return promo array, but got error while test")
	}

	actualUserId := promos[0].UserID
	var user models.User
	if err := models.Dbm.SelectOne(&user, "select * from users where id=?", userId); err != nil {
		t.Errorf("can't select user")
	}

	if userId != actualUserId {
		t.Errorf("relation mismatch, expected user_id=%d but got user_id=%d", userId, actualUserId)
	}
}

func TestUpdatePromo(t *testing.T) {
	var promo models.Promo
	if err := models.Dbm.SelectOne(&promo, "select * from promos where id=?", 1); err != nil {
		t.Errorf("can't select data, expected return promo, but got error while test")
	}

	newName := "New Name"
	promo.Name = newName
	if err := promo.Update(); err != nil {
		t.Errorf("can't update promo")
	}

	var promoUpdated models.Promo
	if err := models.Dbm.SelectOne(&promoUpdated, "select * from promos where id=?", 1); err != nil {
		t.Errorf("can't select data, expected return promo, but got error while test")
	}

	if promoUpdated.Name != newName {
		t.Errorf("name mismatch, expected name=%s but got name=%s", newName, promoUpdated.Name)
	}
}
