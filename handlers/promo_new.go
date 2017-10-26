package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/99ridho/come-backend/env"
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

	broadcastToNearbyUsers(req.Latitude, req.Longitude, "New Promo!", fmt.Sprintf("New promo nearby to you was discovered!"), map[string]string{
		"promo_name":        req.Name,
		"promo_description": req.Description,
		"promo_address":     req.Address,
	})

	json.NewEncoder(w).Encode(map[string]string{
		"message": "add promo success",
	})
}

func broadcastToNearbyUsers(promoLatitude, promoLongitude float64, title, message string, data map[string]string) error {
	url := "https://onesignal.com/api/v1/notifications"

	query := `
	SELECT fcm_token FROM users u
	WHERE (
		6371 *
		acos(cos(radians(?)) * 
		cos(radians(u.current_latitude)) * 
		cos(radians(u.current_longitude) - 
		radians(?)) + 
		sin(radians(?)) * 
		sin(radians(u.current_latitude)))
	) < 5
	`

	var usersPlayerIds []string
	if _, err := models.Dbm.Select(&usersPlayerIds, query, promoLatitude, promoLongitude, promoLatitude); err != nil {
		return err
	}

	if len(usersPlayerIds) < 1 {
		return errors.NewError("no nearby", 500)
	}

	requestData := map[string]interface{}{
		"app_id": env.Getenv("ONE_SIGNAL_APP_ID", ""),
		"contents": map[string]string{
			"en": message,
		},
		"headings": map[string]string{
			"en": title,
		},
		"data":               data,
		"include_player_ids": usersPlayerIds,
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
