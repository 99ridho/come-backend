package models

type Promo struct {
	ID            int     `db:"id"`
	UserID        int     `db:"user_id"`
	Name          string  `db:"name"`
	Address       string  `db:"address"`
	Latitude      float64 `db:"latitude"`
	Longitude     float64 `db:"longitude"`
	Description   string  `db:"description"`
	StartDate     string  `db:"start_date"`     // yyyy-MM-dd HH:mm
	EndDate       string  `db:"end_date"`       // yyyy-MM-dd HH:mm
	AllowedGender string  `db:"allowed_gender"` // male, female or both
	MaxSlot       int     `db:"max_slot"`
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
	if err := Dbm.SelectOne(user, "select * from user where id=?", p.UserID); err != nil {
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
