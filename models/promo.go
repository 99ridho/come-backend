package models

type Promo struct {
	ID            int     `db:"id" json:"id"`
	UserID        int     `db:"user_id" json:"user_id"`
	Name          string  `db:"name" json:"name"`
	Address       string  `db:"address" json:"address"`
	Latitude      float64 `db:"latitude" json:"latitude"`
	Longitude     float64 `db:"longitude" json:"longitude"`
	Description   string  `db:"description" json:"description"`
	StartDate     string  `db:"start_date" json:"start_date"`         // yyyy-MM-dd HH:mm
	EndDate       string  `db:"end_date" json:"end_date"`             // yyyy-MM-dd HH:mm
	AllowedGender string  `db:"allowed_gender" json:"allowed_gender"` // male, female or both
	MaxSlot       int     `db:"max_slot" json:"max_slot"`
}

func NewPromo(userId int, name string, address string, latitude float64, longitude float64, description string, startDate string, endDate string, allowedGender string, maxSlot int) (*Promo, error) {

	newPromo := &Promo{
		UserID:        userId,
		Name:          name,
		Address:       address,
		Latitude:      latitude,
		Longitude:     longitude,
		Description:   description,
		StartDate:     startDate,
		EndDate:       endDate,
		AllowedGender: allowedGender,
		MaxSlot:       maxSlot,
	}

	if err := Dbm.Insert(newPromo); err != nil {
		return nil, err
	}

	return newPromo, nil
}

func (p *Promo) User() (*User, error) {
	var user *User
	if err := Dbm.SelectOne(user, "select * from users where id=?", p.UserID); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *Promo) Update() error {
	if _, err := Dbm.Update(p); err != nil {
		return err
	}

	return nil
}

func (p *Promo) Delete() error {
	if _, err := Dbm.Delete(p); err != nil {
		return err
	}

	return nil
}
