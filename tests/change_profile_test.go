package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"goji.io/pat"

	"goji.io"

	"github.com/99ridho/come-backend/handlers"
	"github.com/99ridho/come-backend/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestChangeProfile(t *testing.T) {
	strReader := strings.NewReader(`
		{
			"username":"johndoe",
			"full_name":"john doe",
			"gender":"male"
		}`)

	req := httptest.NewRequest("POST", "/change_profile", strReader)

	// login o sek, bar ngunu delehen token e ndek kene cur.
	req.Header.Add("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDc2MDMwMjcsInVzZXJfaWQiOjF9.UWtIhU8PF7AToKSrk_xoIua6HTi35meX_10M0INq_Rs")
	rr := httptest.NewRecorder()
	mux := goji.NewMux()
	mux.Use(middlewares.VerifyToken)
	mux.HandleFunc(pat.Post("/change_profile"), handlers.ChangeProfile)
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	jsonExpected := string(`{
		"message": "profile updated"
	}`)

	assert.JSONEq(t, jsonExpected, rr.Body.String(), "Response different")
}
