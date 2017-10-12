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

	userRoute := goji.SubMux()
	userRoute.Use(middlewares.VerifyToken)
	userRoute.HandleFunc(pat.Post("/change_password"), handlers.ChangePassword)
	userRoute.HandleFunc(pat.Post("/change_profile"), handlers.ChangeProfile)

	promoRoute := goji.SubMux()
	promoRoute.Use(middlewares.VerifyToken)
	promoRoute.HandleFunc(pat.Post("/new"), handlers.NewPromo)
	promoRoute.HandleFunc(pat.Get("/my_promos"), handlers.FetchMyPromos)
	promoRoute.HandleFunc(pat.Post("/nearby_promos"), handlers.FetchNearbyPromosFromMyLocation)
	promoRoute.HandleFunc(pat.Get("/:id"), handlers.FetchPromoById)
	promoRoute.HandleFunc(pat.Delete("/:id"), handlers.DeleteMyPromoById)
	promoRoute.HandleFunc(pat.Post("/:id"), handlers.ChangeMyPromoById)
	promoRoute.HandleFunc(pat.Post("/:id/join"), handlers.JoinPromo)

	rootRoute.Handle(pat.New("/user/*"), userRoute)
	rootRoute.Handle(pat.New("/promo/*"), promoRoute)
	return rootRoute
}
