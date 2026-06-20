package api_functions

import (
	"context"
	"encoding/json"
	"log"

	"github.com/naki/dispatch-service/conf"
	"github.com/naki/dispatch-service/models"
	"github.com/segmentio/kafka-go"
)

func StartKafkaConsumers() {
	go consumeBookingCreated()
}

func consumeBookingCreated() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{conf.AppConfig.KafkaBroker},
		Topic:    "booking.created",
		GroupID:  "dispatch-service",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	log.Println("listening on kafka topic: booking.created")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading from topic booking.created: %v", err)
			continue
		}

		var event models.BookingEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshal booking event: %v", err)
			continue
		}

		log.Printf("received booking.created: id=%s type=%s service=%s",
			event.BookingID, event.BookingType, event.ServiceType)

		go handleBookingCreated(event)
	}
}

func handleBookingCreated(event models.BookingEvent) {
	result, err := FindBestNurse(event)
	if err != nil {
		log.Printf("dispatch failed for booking %s: %v", event.BookingID, err)
		LogDispatchFailure(event, err.Error())
		return
	}

	log.Printf("dispatch matched booking %s -> nurse %s (distance=%.2fkm score=%.4f)",
		event.BookingID, result.NurseID, result.Distance, result.MatchScore)

	if err := SetNurseInSession(result.NurseID, true); err != nil {
		log.Printf("warning: could not set nurse %s in-session: %v", result.NurseID, err)
	}

	PublishNurseMatched(event, result)
}
