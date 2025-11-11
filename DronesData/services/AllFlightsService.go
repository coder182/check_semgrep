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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// to make date formate according to API date formate
const apiTimeLayout = "2006-01-02T15:04:05.000Z"

func StoreAllFlightsData() error {

	var startDate, endDate time.Time
	//  Prepare base URL
	baseURL := os.Getenv("PAGINATION_API_URL")
	dateLimit := os.Getenv("DRON_RECORD_TIME_LIMIT")

	if baseURL == "" {
		log.Fatal("Pangination API URL doesn't exists in env file")
	}
	log.Printf("Pagination API URL loaded from env file")
	log.Printf("Checking last inserted record in Flights_data collection")
	collection := database.AllFlightsDataCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//  Step 1: Check if any data already exists
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to count documents: %v", err)
	}

	var finalURL string

	if count == 0 {
		// if datelimit is not set from env
		if dateLimit == "" {
			finalURL = baseURL
			log.Println("No data found in DB. Fetching all data from API")
		} else {
			//Reading time limit from env
			startDate, err = parseDateLimit(dateLimit)
			if err != nil {
				log.Printf("Failed to parse DRON_RECORD_TIME_LIMIT: %v. Fetching all data", err)
				finalURL = baseURL
			} else {
				endDate = time.Now().UTC()
				finalURL = fmt.Sprintf("%s?startDate=%s&endDate=%s",
					baseURL,
					startDate.Format(apiTimeLayout),
					endDate.Format(apiTimeLayout),
				)
				log.Printf("DB empty. Fetching data from %s to %s",
					startDate.Format(apiTimeLayout), endDate.Format(apiTimeLayout))
			}
		}
	} else {
		log.Printf("Found record in Flights_data collection ")
		log.Printf("Getting last inserted record time")

		//Read latest end_time from DB
		opts := options.FindOne().
			SetSort(bson.D{{"end_time", -1}}).
			SetProjection(bson.M{"end_time": 1})

		var latest model.DronesTrack
		if err := collection.FindOne(ctx, bson.M{}, opts).Decode(&latest); err != nil {
			return fmt.Errorf("failed to fetch latest end_time: %v", err)
		}

		startDate = latest.EndTime.Time
		endDate = time.Now().UTC()

		if dateLimit != "" {
			log.Printf("WARNING: DRON_RECORD_TIME_LIMIT is provided but IGNORED for subsequent runs. Using latest DB timestamp instead.")
		}

		log.Printf("Start date from DB: %s", startDate.Format(apiTimeLayout))
		log.Printf("End date (current time): %s", endDate.Format(apiTimeLayout))

		// Build API URL with startDate & endDate
		endDate = time.Now().UTC()
		finalURL = fmt.Sprintf("%s?startDate=%s&endDate=%s",
			baseURL,
			startDate.Format(apiTimeLayout),
			endDate.Format(apiTimeLayout),
		)
		log.Printf("Fetching new data from %s to %s", startDate.Format(apiTimeLayout), endDate.Format(apiTimeLayout))

	}

	//Make API Call
	req, err := api.NewAuthenticatedRequest("GET", finalURL)
	if err != nil {
		return fmt.Errorf("failed to create authenticated request: %v", err)
	}

	resp, err := api.AuthClient.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Check status after reading body
	if resp.StatusCode != http.StatusOK {
		log.Printf("Status code %d", resp.StatusCode)
		return fmt.Errorf("API request failed: %d - %s", resp.StatusCode, string(body))
	}

	// Parse the API response
	var apiResponse model.AllFlightsResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// Step 3: Insert new records to MongoDB
	if len(apiResponse.Data) == 0 {
		log.Println("No new flights data found.")
		return nil
	}

	var docs []interface{}
	for _, flight := range apiResponse.Data {
		docs = append(docs, flight)
	}

	result, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to insert flights data: %v", err)
	}

	log.Printf("Inserted %d flight records", len(result.InsertedIDs))
	log.Printf("######################################################")
	return nil

}

// Helper function for date parsing
func parseDateLimit(dateLimit string) (time.Time, error) {
	// Try RFC3339 first (with timezone)
	if parseDate, err := time.Parse(time.RFC3339, dateLimit); err == nil {
		return parseDate, nil
	}

	// Try without timezone (assume UTC)
	return time.Parse("2006-01-02T15:04:05", dateLimit)
}
