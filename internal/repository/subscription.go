package repository

import (
	"subscriptions-service/internal/model"
	"time"

	"github.com/google/uuid"
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
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query,
		s.ServiceName,
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

func (r *SubscriptionRepo) GetByID(id int) (model.Subscription, error) {
	var sub model.Subscription

	query := `SELECT * FROM subscriptions WHERE id=$1`

	err := r.db.Get(&sub, query, id)
	return sub, err
}

func (r *SubscriptionRepo) Update(sub model.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name=$1, price=$2, user_id=$3, start_date=$4, end_date=$5
		WHERE id=$6
	`

	_, err := r.db.Exec(query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)

	return err
}

func (r *SubscriptionRepo) Delete(id int) error {
	query := `DELETE FROM subscriptions WHERE id=$1`

	_, err := r.db.Exec(query, id)
	return err
}

func (r *SubscriptionRepo) GetTotalSum(userID uuid.UUID, serviceName string, from, to *time.Time) (int, error) {
	var sum int

	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE ($1::uuid IS NULL OR user_id = $1)
		  AND ($2 = '' OR service_name = $2)
		  AND ($3::date IS NULL OR start_date >= $3)
		  AND ($4::date IS NULL OR start_date <= $4)
	`

	err := r.db.Get(&sum, query, userID, serviceName, from, to)
	return sum, err
}