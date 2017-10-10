package main

import (
	"github.com/99ridho/come-backend/models"
)

func main() {
	if err := models.CreateTables(); err != nil {
		panic(err)
	}

	models.NewUser("test", "test@mail.com", "123321", "Test Account", "male", "8912830812387192837")
	models.NewPromo(1, "McD Gratis Big Mac", "Jalan MT. Haryono Malang", -7.23881238, 122.99812, "McD Gratis Big Mac", "2017-10-05 12:00", "2017-10-05 12:00", "both", 5)
}
