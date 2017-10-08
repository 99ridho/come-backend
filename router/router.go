package router

import (
	"goji.io"
	"goji.io/pat"

	"github.com/99ridho/come-backend/handlers"
	"github.com/99ridho/come-backend/middlewares"
)

func NewRouter() *goji.Mux {
	rootRoute := goji.NewMux()
	rootRoute.HandleFunc(pat.Post("/login"), handlers.Login)
	rootRoute.HandleFunc(pat.Post("/register"), handlers.Register)

	routeWithAuth := goji.SubMux()
	routeWithAuth.Use(middlewares.VerifyToken)
	// router with auth added here...

	rootRoute.Handle(pat.New("/*"), routeWithAuth)
	return rootRoute
}
