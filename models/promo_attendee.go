package models

type PromoAttendee struct {
	ID      int `db:"id"`
	UserID  int `db:"user_id"`
	PromoID int `db:"promo_id"`
}

func NewPromoAttendee(userId int, promoId int) (*PromoAttendee, error) {
	newPromoAttendee := &PromoAttendee{
		UserID:  userId,
		PromoID: promoId,
	}

	if err := Dbm.Insert(newPromoAttendee); err != nil {
		return nil, err
	}

	return newPromoAttendee, nil
}

func (p *PromoAttendee) User() (*User, error) {
	var user User
	if err := Dbm.SelectOne(&user, "select * from users where id=?", p.UserID); err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *PromoAttendee) Promo() (*Promo, error) {
	var promo Promo
	if err := Dbm.SelectOne(&promo, "select * from promos where id=?", p.PromoID); err != nil {
		return nil, err
	}

	return &promo, nil
}
