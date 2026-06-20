package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/naki/dispatch-service/functions/api_functions"
	"github.com/naki/dispatch-service/models"
)

func GoOnline(c *gin.Context) {
	nurseID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "invalid user id"})
		return
	}

	var body struct {
		Latitude  float64  `json:"latitude"`
		Longitude float64  `json:"longitude"`
		Services  []string `json:"services"`
		Rating    float64  `json:"rating"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 400, "message": "latitude and longitude are required"})
		return
	}

	if err := api_functions.SetNurseOnline(nurseID, body.Latitude, body.Longitude, body.Services, body.Rating); err != nil {
		log.Printf("failed to set nurse %s online: %v", nurseID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": 500, "message": "failed to go online"})
		return
	}

	log.Printf("nurse %s went online at (%.6f, %.6f)", nurseID, body.Latitude, body.Longitude)

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "you are now online and available for bookings",
	})
}

func GoOffline(c *gin.Context) {
	nurseID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "invalid user id"})
		return
	}

	if err := api_functions.SetNurseOffline(nurseID); err != nil {
		log.Printf("failed to set nurse %s offline: %v", nurseID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": 500, "message": "failed to go offline"})
		return
	}

	log.Printf("nurse %s went offline", nurseID)

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "you are now offline",
	})
}

func UpdateLocation(c *gin.Context) {
	nurseID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "invalid user id"})
		return
	}

	var body struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 400, "message": "latitude and longitude required"})
		return
	}

	if err := api_functions.UpdateNurseLocation(nurseID, body.Latitude, body.Longitude); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "location updated",
	})
}

func GetAvailableNurses(c *gin.Context) {
	nurses, err := api_functions.GetAllAvailableNurses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 500, "message": err.Error()})
		return
	}

	if nurses == nil {
		nurses = []models.NurseAvailability{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "available nurses retrieved",
		"data":    nurses,
		"count":   len(nurses),
	})
}

func ManualDispatch(c *gin.Context) {
	var body models.BookingEvent
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 400, "message": "invalid booking data"})
		return
	}

	result, err := api_functions.FindBestNurse(body)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": 404, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "nurse matched",
		"data":    result,
	})
}

func GetDispatchHistory(c *gin.Context) {
	bookingID, err := uuid.Parse(c.Param("booking_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 400, "message": "invalid booking id"})
		return
	}

	logs, err := api_functions.GetDispatchLogs(bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "dispatch history retrieved",
		"data":    logs,
	})
}

func GetRecentDispatches(c *gin.Context) {
	logs, err := api_functions.GetRecentDispatches(50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "recent dispatches retrieved",
		"data":    logs,
	})
}
