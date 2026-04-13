package repository

import (
	"subscriptions-service/internal/model"

	"github.com/jmoiron/sqlx"
)

type SubscriptionRepo struct {
	db *sqlx.DB
}

func NewSubscriptionRepo(db *sqlx.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(s model.Subscription) error {
	query := `
		INSERT INTO subscriptions (name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query,
		s.Name,
		s.Price,
		s.UserID,
		s.StartDate,
		s.EndDate,
	)

	return err
}

func (r *SubscriptionRepo) GetAll() ([]model.Subscription, error) {
	var subs []model.Subscription

	query := `SELECT * FROM subscriptions`

	err := r.db.Select(&subs, query)
	return subs, err
}