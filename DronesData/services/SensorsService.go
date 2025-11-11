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

// FetchAndStoreDrones fetches drone data and stores in MongoDB
func FetchAndStoreSensors() error {

	// Build API URL from environment
	baseURL := os.Getenv("SENSOR_API_URL")

	if baseURL == "" {
		log.Fatal("Sensors API URL doesn't exists in the properties file")
	}
	log.Printf("Sensors API URL loaded from env file")
	// Create authenticated request
	req, err := api.NewAuthenticatedRequest("GET", baseURL)
	if err != nil {
		return err
	}

	// Execute the request
	resp, err := api.AuthClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed: %s - %s", resp.Status, string(body))
	}
	log.Printf("Sensors API URL OATH successful")
	// 2. Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 3. Parse JSON response
	var apiResponse []struct {
		Serial   string `json:"serial"`
		Location struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"location"`
		Type             *string  `json:"type"`
		Name             string   `json:"name"`
		HeightAGL        *float64 `json:"heightAGL"`
		HeightGeodetic   float64  `json:"heightGeodetic"`
		LastOnline       string   `json:"lastOnline"` // String representation
		Pressure         float64  `json:"pressure"`
		Temperature      float64  `json:"temperature"`
		IPAddress        *string  `json:"ipAddress"`
		WarningRadius    *float64 `json:"warningRadius"`
		AlertRadius      *float64 `json:"alertRadius"`
		Online           bool     `json:"online"`
		State            int      `json:"state"`
		NoiseFloorValue  float64  `json:"noiseFloorValue"`
		NoiseFloorStatus string   `json:"noiseFloorStatus"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return err
	}

	// 4. Convert to our model
	var drones []model.SensorsData
	for _, item := range apiResponse {
		// Parse LastOnline string to time.Time
		lastOnline, err := time.Parse("2006-01-02 15:04:05", item.LastOnline)
		if err != nil {
			log.Printf("Failed to parse time for drone %s: %v", item.Serial, err)
			continue
		}

		drones = append(drones, model.SensorsData{
			Serial:           item.Serial,
			Location:         model.Location(item.Location),
			Type:             item.Type,
			Name:             item.Name,
			HeightAGL:        item.HeightAGL,
			HeightGeodetic:   item.HeightGeodetic,
			LastOnline:       lastOnline,
			Pressure:         item.Pressure,
			Temperature:      item.Temperature,
			IPAddress:        item.IPAddress,
			WarningRadius:    item.WarningRadius,
			AlertRadius:      item.AlertRadius,
			Online:           item.Online,
			State:            item.State,
			NoiseFloorValue:  item.NoiseFloorValue,
			NoiseFloorStatus: item.NoiseFloorStatus,
		})
	}

	// 5. Insert into MongoDB
	collection := database.SensorsDataCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert to interface slice for bulk insert
	var docs []interface{}
	for _, drone := range drones {
		docs = append(docs, drone)
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	log.Printf("Inserted %d drone records", len(result.InsertedIDs))
	log.Printf("###################################################")
	return nil
}
