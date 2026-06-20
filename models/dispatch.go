package models

import (
	"time"

	"github.com/google/uuid"
)

type NurseAvailability struct {
	NurseID    uuid.UUID `json:"nurse_id"    redis:"nurse_id"`
	IsOnline   bool      `json:"is_online"   redis:"is_online"`
	Latitude   float64   `json:"latitude"    redis:"latitude"`
	Longitude  float64   `json:"longitude"   redis:"longitude"`
	Services   []string  `json:"services"    redis:"services"`
	Rating     float64   `json:"rating"      redis:"rating"`
	Verified   bool      `json:"verified"    redis:"verified"`
	InSession  bool      `json:"in_session"  redis:"in_session"`
	UpdatedAt  time.Time `json:"updated_at"  redis:"updated_at"`
}

type DispatchLog struct {
	ID          uuid.UUID  `db:"id"           json:"id"`
	BookingID   uuid.UUID  `db:"booking_id"   json:"booking_id"`
	NurseID     *uuid.UUID `db:"nurse_id"     json:"nurse_id"`
	Status      string     `db:"status"       json:"status"`
	Reason      string     `db:"reason"       json:"reason"`
	Distance    float64    `db:"distance"     json:"distance"`
	MatchScore  float64    `db:"match_score"  json:"match_score"`
	BookingType string     `db:"booking_type" json:"booking_type"`
	CreatedAt   time.Time  `db:"created_at"   json:"created_at"`
}

type MatchResult struct {
	NurseID    uuid.UUID `json:"nurse_id"`
	Distance   float64   `json:"distance_km"`
	Rating     float64   `json:"rating"`
	MatchScore float64   `json:"match_score"`
}

type BookingEvent struct {
	BookingID     string `json:"booking_id"`
	CustomerID    string `json:"customer_id"`
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	CustomerEmail string `json:"customer_email"`
	ServiceType   string `json:"service_type"`
	BookingType   string `json:"booking_type"`
	ScheduledAt   string `json:"scheduled_at"`
	Address       string `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Status        string `json:"status"`
}
