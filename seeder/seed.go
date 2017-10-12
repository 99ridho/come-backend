package main

import (
	"github.com/99ridho/come-backend/models"
)

func main() {
	if err := models.CreateTables(); err != nil {
		panic(err)
	}

	models.NewUser("test", "test@mail.com", "123321", "Test Account", "male", "8912830812387192837")                                                                        // 1
	models.NewUser("test1", "test1@mail.com", "123321", "Test Account 1", "female", "8912830812387192837")                                                                  // 2
	models.NewPromo(1, "McD Gratis Big Mac", "Jalan MT. Haryono Malang", -7.23881238, 122.99812, "McD Gratis Big Mac", "2017-10-05 12:00", "2017-10-05 12:00", "both", 5)   // 1
	models.NewPromo(1, "McD Gratis Mc Float", "Jalan MT. Haryono Malang", -7.23881238, 122.99812, "McD Gratis Mc Float", "2017-10-05 12:00", "2017-10-05 12:00", "both", 5) // 2
	// models.NewPromoAttendee(2, 1)
	// models.NewPromoAttendee(2, 2)
}
