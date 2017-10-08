package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99ridho/come-backend/handlers"
	"github.com/stretchr/testify/assert"
)

func TestLoginWithUnknownAccount(t *testing.T) {
	req := httptest.NewRequest("POST", "/login", strings.NewReader(`{"email":"aaa@gmail.com","pasword":"1111"}`))
	rr := httptest.NewRecorder()
	http.HandlerFunc(handlers.Login).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	jsonExpected := string(`{"message":"user not found","status":"failed","token":""}`)

	print(rr.Body.String())
	assert.JSONEq(t, jsonExpected, rr.Body.String(), "Response different")
}
