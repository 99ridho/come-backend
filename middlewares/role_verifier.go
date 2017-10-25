package middlewares

import (
	"net/http"

	"github.com/99ridho/come-backend/errors"
	"github.com/99ridho/come-backend/models"
)

func VerifyRoleAdmin(h http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user_id").(int)
		var user models.User

		if err := models.Dbm.SelectOne(&user, "select * from users where id=?", userId); err != nil {
			errors.NewErrorWithStatusCode(http.StatusInternalServerError).WriteTo(w)
			return
		}

		if user.Role != "admin" {
			errors.NewError("only admin can execute this endpoint", http.StatusUnauthorized).WriteTo(w)
			return
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(middleware)
}
