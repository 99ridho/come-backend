package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"goji.io/pat"

	"github.com/99ridho/come-backend/env"
	"github.com/99ridho/come-backend/errors"
	"github.com/99ridho/come-backend/models"
)

func JoinPromo(w http.ResponseWriter, r *http.Request) {
	promoId, _ := strconv.ParseInt(pat.Param(r, "id"), 10, 0)

	// 1. select user
	userId := r.Context().Value("user_id").(int)
	var userWantToJoin models.User
	if err := models.Dbm.SelectOne(&userWantToJoin, "select * from users where id=?", userId); err != nil {
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

	// 3. check if already registered
	attendeeCountWithUserIdWantToJoin, countError := models.Dbm.SelectInt("select count(*) from promo_attendees where promo_id=? and user_id=?", promoId, userId)
	if countError != nil {
		errors.NewErrorWithStatusCode(http.StatusInternalServerError).WriteTo(w)
		return
	}
	if int(attendeeCountWithUserIdWantToJoin) > 0 {
		errors.NewError("already joined", http.StatusInternalServerError).WriteTo(w)
		return
	}

	// 4. check allowed gender
	if promo.AllowedGender != "both" {
		if promo.AllowedGender != userWantToJoin.Gender {
			errors.NewError("can't join promo, gender mismatch!", http.StatusInternalServerError).WriteTo(w)
			return
		}
	}

	// 5. check if promo is full slot
	if int(attendeeCount) >= promo.MaxSlot {
		errors.NewError("can't join promo, slot full!", http.StatusInternalServerError).WriteTo(w)
		return
	}

	// 6. check if inserting to DB encountered an error
	if _, err := models.NewPromoAttendee(userId, int(promoId)); err != nil {
		errors.NewError("can't join promo", http.StatusInternalServerError).WriteTo(w)
		return
	}

	// 7. TODO - notify promo owner..
	promoOwner, _ := promo.User()
	println("Promo owner : ", promoOwner.FcmToken)
	sendNotificationToPromoOwner(promoOwner, "Someone Joined", fmt.Sprintf("%s joined your promo named %s!", userWantToJoin.FullName, promo.Name))

	json.NewEncoder(w).Encode(map[string]string{
		"message": "joining promo success",
	})
}

func sendNotificationToPromoOwner(owner *models.User, title, message string) error {
	url := "https://onesignal.com/api/v1/notifications"

	requestData := map[string]interface{}{
		"app_id": env.Getenv("ONE_SIGNAL_APP_ID", ""),
		"contents": map[string]string{
			"en": message,
		},
		"headings": map[string]string{
			"en": title,
		},
		"include_player_ids": []string{owner.FcmToken},
	}

	encodedRequest, _ := json.Marshal(requestData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(encodedRequest))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Basic "+env.Getenv("ONE_SIGNAL_REST_API_KEY", ""))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}
