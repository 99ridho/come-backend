package main

import (
	"github.com/99ridho/come-backend/models"
)

func main() {
	if err := models.CreateTables(); err != nil {
		panic(err)
	}

	models.NewUser(1, "test", "test@mail.com", "123321", "Test Account", "male", "8912830812387192837")
}
