package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99ridho/come-backend/handlers"
	"github.com/99ridho/come-backend/models"
	"github.com/stretchr/testify/assert"
)

func TestRegisterNewAccount(t *testing.T) {
	strReader := strings.NewReader(`
		{
			"username":"johndoe",
			"email":"balala@gmail.com",
			"password":"password",
			"full_name":"john doe",
			"gender":"male",
			"fcm_token":"tokentokentoken"	
		}`)

	countBefore, _ := models.Dbm.SelectInt("select count(*) from users")
	req := httptest.NewRequest("POST", "/register", strReader)
	rr := httptest.NewRecorder()
	http.HandlerFunc(handlers.Register).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	countAfter, _ := models.Dbm.SelectInt("select count(*) from users")

	assert.Equal(t, countBefore+1, countAfter, "Count before and after differs")
}

func TestRegisterWithExistingAccount(t *testing.T) {
	strReader := strings.NewReader(`
		{
			"username":"johndoe",
			"email":"balala@gmail.com",
			"password":"password",
			"full_name":"john doe",
			"gender":"male",
			"fcm_token":"tokentokentoken"	
		}`)

	countBefore, _ := models.Dbm.SelectInt("select count(*) from users")
	req := httptest.NewRequest("POST", "/register", strReader)
	rr := httptest.NewRecorder()
	http.HandlerFunc(handlers.Register).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	countAfter, _ := models.Dbm.SelectInt("select count(*) from users")

	assert.Equal(t, countBefore, countAfter, "Count before and after differs")
}
