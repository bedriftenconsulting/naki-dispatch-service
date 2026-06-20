package api_functions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/naki/dispatch-service/conf"
	"github.com/naki/dispatch-service/models"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     conf.AppConfig.RedisAddr,
		Password: conf.AppConfig.RedisPass,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("WARNING: redis not available: %v (dispatch will use DB fallback)", err)
	} else {
		log.Println("redis connected successfully")
	}
}

func nurseKey(nurseID uuid.UUID) string {
	return fmt.Sprintf("nurse:available:%s", nurseID.String())
}

func SetNurseAvailability(nurse models.NurseAvailability) error {
	ctx := context.Background()

	data, err := json.Marshal(nurse)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, nurseKey(nurse.NurseID), data, 30*time.Minute).Err()
}

func GetNurseAvailability(nurseID uuid.UUID) (*models.NurseAvailability, error) {
	ctx := context.Background()

	data, err := rdb.Get(ctx, nurseKey(nurseID)).Bytes()
	if err != nil {
		return nil, err
	}

	var nurse models.NurseAvailability
	if err := json.Unmarshal(data, &nurse); err != nil {
		return nil, err
	}

	return &nurse, nil
}

func RemoveNurseAvailability(nurseID uuid.UUID) error {
	ctx := context.Background()
	return rdb.Del(ctx, nurseKey(nurseID)).Err()
}

func GetAllAvailableNurses() ([]models.NurseAvailability, error) {
	ctx := context.Background()

	keys, err := rdb.Keys(ctx, "nurse:available:*").Result()
	if err != nil {
		return nil, err
	}

	var nurses []models.NurseAvailability

	for _, key := range keys {
		data, err := rdb.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var nurse models.NurseAvailability
		if err := json.Unmarshal(data, &nurse); err != nil {
			continue
		}

		if nurse.IsOnline && !nurse.InSession && nurse.Verified {
			nurses = append(nurses, nurse)
		}
	}

	return nurses, nil
}

func UpdateNurseLocation(nurseID uuid.UUID, lat, lng float64) error {
	nurse, err := GetNurseAvailability(nurseID)
	if err != nil {
		return fmt.Errorf("nurse not found in availability pool")
	}

	nurse.Latitude = lat
	nurse.Longitude = lng
	nurse.UpdatedAt = time.Now()

	return SetNurseAvailability(*nurse)
}

func SetNurseInSession(nurseID uuid.UUID, inSession bool) error {
	nurse, err := GetNurseAvailability(nurseID)
	if err != nil {
		return err
	}

	nurse.InSession = inSession
	nurse.UpdatedAt = time.Now()

	return SetNurseAvailability(*nurse)
}
