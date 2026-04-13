package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int        `json:"id"`
	Name 		string     `json:"name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}