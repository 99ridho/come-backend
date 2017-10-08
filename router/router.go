package router

import (
	"net/http"

	"goji.io"
	"goji.io/pat"

	"encoding/json"

	"github.com/99ridho/come-backend/handlers"
)

func NewRouter() *goji.Mux {
	router := goji.NewMux()

	router.HandleFunc(pat.Get("/hello"), func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Hello",
		})
	})

	router.HandleFunc(pat.Post("/login"), handlers.Login)

	return router
}
