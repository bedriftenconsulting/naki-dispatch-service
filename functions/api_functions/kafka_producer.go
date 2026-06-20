package api_functions

import (
	"context"
	"encoding/json"
	"log"

	"github.com/naki/dispatch-service/conf"
	"github.com/naki/dispatch-service/models"
	"github.com/segmentio/kafka-go"
)

type NurseMatchedEvent struct {
	BookingID     string  `json:"booking_id"`
	CustomerID    string  `json:"customer_id"`
	NurseID       string  `json:"nurse_id"`
	CustomerName  string  `json:"customer_name"`
	CustomerPhone string  `json:"customer_phone"`
	ServiceType   string  `json:"service_type"`
	ScheduledAt   string  `json:"scheduled_at"`
	Address       string  `json:"address"`
	Distance      float64 `json:"distance_km"`
	MatchScore    float64 `json:"match_score"`
}

func PublishNurseMatched(booking models.BookingEvent, result *models.MatchResult) {
	event := NurseMatchedEvent{
		BookingID:     booking.BookingID,
		CustomerID:    booking.CustomerID,
		NurseID:       result.NurseID.String(),
		CustomerName:  booking.CustomerName,
		CustomerPhone: booking.CustomerPhone,
		ServiceType:   booking.ServiceType,
		ScheduledAt:   booking.ScheduledAt,
		Address:       booking.Address,
		Distance:      result.Distance,
		MatchScore:    result.MatchScore,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal nurse.matched event: %v", err)
		return
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{conf.AppConfig.KafkaBroker},
		Topic:   "nurse.matched",
	})
	defer writer.Close()

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(booking.BookingID),
			Value: payload,
		},
	)

	if err != nil {
		log.Printf("failed to publish nurse.matched event: %v", err)
		return
	}

	log.Printf("nurse.matched event published: booking=%s nurse=%s", booking.BookingID, result.NurseID)
}
