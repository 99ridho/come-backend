package main

import (
	"net/http"

	r "github.com/99ridho/come-backend/router"
)

func main() {
	router := r.NewRouter()
	http.ListenAndServe(":8888", router)
}
