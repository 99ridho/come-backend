package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/99ridho/come-backend/env"

	jwt "github.com/dgrijalva/jwt-go"
)

func VerifyToken(h http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		userId, err := checkToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "unauthorized",
				"message": err.Error(),
			})
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", int(userId))
		reqWithContext := r.WithContext(ctx)
		h.ServeHTTP(w, reqWithContext)
	}

	return http.HandlerFunc(middleware)
}

func checkToken(tokenString string) (id float64, err error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secret := env.Getenv("SECRET_KEY", "secret")
		return []byte(secret), nil
	})

	if err != nil {
		return id, errors.New("invalid token : " + err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return id, errors.New("invalid token: claim invalid")
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return id, errors.New("invalid token: token expired")
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return id, errors.New("invalid token: id isn't valid int number")
	}

	return userId, nil
}
