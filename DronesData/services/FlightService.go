package services

import (
	"DronesData/api"
	"DronesData/database"
	model "DronesData/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func StoreFlightLocation(serialID string) error {

	baseURL := os.Getenv("LOCATION_API_URL")
	if baseURL == "" {
		log.Fatal("Drones Location API doesn't exists in the properties file")
	}
	log.Printf("Location API URL Loaded from env file")
	// Make authenticated request
	req, err := api.NewAuthenticatedRequest("GET", baseURL+serialID)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := api.AuthClient.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %v", err)
	}
	log.Printf("Location API URL OATH successful")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned status %d for track ID %s", resp.StatusCode, serialID)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	// Parse as array
	var flights []model.FlightLocation
	if err := json.Unmarshal(body, &flights); err != nil {
		return fmt.Errorf("failed to parse flight data: %v", err)
	}
	log.Printf("Received %d flight records for track ID: %s", len(flights), serialID)

	// Insert into MongoDB
	collection := database.FlightLocationCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert to slice of interfaces
	var docs []interface{}
	for _, flight := range flights {
		docs = append(docs, flight)
	}

	if len(docs) == 0 {
		return fmt.Errorf("no flight data received")
	}

	insertResult, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to insert flights: %v", err)
	}
	log.Printf("Successfully inserted %d records for track ID: %s", len(insertResult.InsertedIDs), serialID)
	log.Printf("#####################")

	return nil
}
