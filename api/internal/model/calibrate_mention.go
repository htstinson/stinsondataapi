package model

import (
	"time"

	"github.com/google/uuid"
)

type CalibrateMention struct {
	ID                uuid.UUID `json:"id" db:"id"`
	CreatedAt         time.Time `json:"created_at"`
	ModifiedAt        time.Time `json:"modified_at"`
	CalibrateResultID uuid.UUID `json:"calibrate_result_id"`
	Rating            int       `json:"rating"`
	SubscriberID      uuid.UUID `json:"subscriber_id"`
	Title             *string   `json:"title"`
	Body              *string   `json:"body"`
	Author            *string   `json:"author"`
	Location          *string   `json:"location"`
	Headline          *string   `json:"headline"`
	RatingDate        time.Time `json:"search_time"`
}
