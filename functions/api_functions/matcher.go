package api_functions

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/naki/dispatch-service/database"
	"github.com/naki/dispatch-service/models"
)

const (
	maxDistanceKm     = 20.0
	distanceWeight    = 0.4
	ratingWeight      = 0.3
	emergencyBonus    = 0.3
	earthRadiusKm     = 6371.0
)

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

func FindBestNurse(booking models.BookingEvent) (*models.MatchResult, error) {
	nurses, err := GetAllAvailableNurses()
	if err != nil {
		return nil, fmt.Errorf("failed to get available nurses: %w", err)
	}

	if len(nurses) == 0 {
		return nil, fmt.Errorf("no available nurses")
	}

	var candidates []models.MatchResult

	for _, nurse := range nurses {
		if !nurse.Verified {
			continue
		}

		if nurse.InSession {
			continue
		}

		if booking.ServiceType != "" && len(nurse.Services) > 0 {
			matched := false
			for _, s := range nurse.Services {
				if s == booking.ServiceType {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		dist := haversine(booking.Latitude, booking.Longitude, nurse.Latitude, nurse.Longitude)

		if dist > maxDistanceKm {
			continue
		}

		distScore := 1.0 - (dist / maxDistanceKm)
		ratingScore := nurse.Rating / 5.0

		score := distScore*distanceWeight + ratingScore*ratingWeight

		if booking.BookingType == "emergency" {
			score += emergencyBonus
			distScore *= 1.5
		}

		candidates = append(candidates, models.MatchResult{
			NurseID:    nurse.NurseID,
			Distance:   math.Round(dist*100) / 100,
			Rating:     nurse.Rating,
			MatchScore: math.Round(score*10000) / 10000,
		})
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no nurses within range for this booking")
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].MatchScore > candidates[j].MatchScore
	})

	best := candidates[0]

	logDispatch(booking, &best)

	return &best, nil
}

func logDispatch(booking models.BookingEvent, result *models.MatchResult) {
	bookingID, err := uuid.Parse(booking.BookingID)
	if err != nil {
		return
	}

	query := `
		INSERT INTO dispatch_logs (booking_id, nurse_id, status, reason, distance, match_score, booking_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	status := "matched"
	reason := fmt.Sprintf("auto-matched: distance=%.2fkm rating=%.1f score=%.4f",
		result.Distance, result.Rating, result.MatchScore)

	_, err = database.DB.Exec(query, bookingID, result.NurseID, status, reason,
		result.Distance, result.MatchScore, booking.BookingType)
	if err != nil {
		log.Printf("failed to log dispatch: %v", err)
	}
}

func LogDispatchFailure(booking models.BookingEvent, failReason string) {
	bookingID, err := uuid.Parse(booking.BookingID)
	if err != nil {
		return
	}

	query := `
		INSERT INTO dispatch_logs (booking_id, status, reason, booking_type)
		VALUES ($1, $2, $3, $4)`

	_, err = database.DB.Exec(query, bookingID, "failed", failReason, booking.BookingType)
	if err != nil {
		log.Printf("failed to log dispatch failure: %v", err)
	}
}

func GetDispatchLogs(bookingID uuid.UUID) ([]models.DispatchLog, error) {
	var logs []models.DispatchLog
	query := `SELECT * FROM dispatch_logs WHERE booking_id = $1 ORDER BY created_at DESC`
	err := database.DB.Select(&logs, query, bookingID)
	return logs, err
}

func GetRecentDispatches(limit int) ([]models.DispatchLog, error) {
	var logs []models.DispatchLog
	query := `SELECT * FROM dispatch_logs ORDER BY created_at DESC LIMIT $1`
	err := database.DB.Select(&logs, query, limit)
	return logs, err
}

func SetNurseOnline(nurseID uuid.UUID, lat, lng float64, services []string, rating float64) error {
	nurse := models.NurseAvailability{
		NurseID:   nurseID,
		IsOnline:  true,
		Latitude:  lat,
		Longitude: lng,
		Services:  services,
		Rating:    rating,
		Verified:  true,
		InSession: false,
		UpdatedAt: time.Now(),
	}

	return SetNurseAvailability(nurse)
}

func SetNurseOffline(nurseID uuid.UUID) error {
	return RemoveNurseAvailability(nurseID)
}
