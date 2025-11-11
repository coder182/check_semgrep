package services

import (
	"DronesData/database"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 2 hours timewindow
const freshWindow = 2 * time.Hour

func ProcessTrackIds() error {
	log.Printf("Reading track ids from track_id collection to process..")
	collection := database.TrackIdsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get limit from environment
	limit := getTrackProcessLimit()
	log.Printf("Read Limit of track ids from env file: %d", limit)

	// Find pending and updating track_ids with applied limit
	filter := bson.M{
		"status": bson.M{"$in": []string{"pending", "updating"}},
	}
	opts := options.Find().SetLimit(int64(limit)).
		SetProjection(bson.M{
			"track_id": 1,
			"end_time": 1,
			"status":   1,
			"_id":      0,
		})

	var items []struct {
		TrackID string    `bson:"track_id"`
		EndTime time.Time `bson:"end_time"`
		Status  string    `bson:"status"`
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return fmt.Errorf("failed to fetch pending track IDs: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return fmt.Errorf("error reading cursor: %w", err)
	}

	log.Printf("Found %d pending and updating track IDs\n", len(items))

	now := time.Now().UTC()

	// Step 2: Loop and process each track ID
	for _, it := range items {

		log.Printf("Processing track ID: %s with endTime: %s", it.TrackID, it.EndTime)

		//empty track id check
		if it.TrackID == "" {
			log.Printf("Skip empty track_id")
			continue
		}

		// if time is zero like January 1, year 1, 00:00:00 UTC.
		if it.EndTime.IsZero() {
			log.Printf("track_id=%s has empty end_time  marking failed", it.TrackID)
			_ = UpdateTrackIdStatus(it.TrackID, "failed")
			continue
		}
		age := now.Sub(it.EndTime.UTC())
		if age >= freshWindow {
			log.Printf("track_id=%s eligible (age=%v)  processing", it.TrackID, age)
			_ = UpdateTrackIdStatus(it.TrackID, "processing")

			err := StoreFlightLocation(it.TrackID)
			if err != nil {
				log.Printf("Failed to store flight location for track ID %s: %v", it.TrackID, err)
				UpdateTrackIdStatus(it.TrackID, "failed")
				continue
			}
			_ = UpdateTrackIdStatus(it.TrackID, "done")
		} else {
			if it.Status != "updating" {
				log.Printf("track_id=%s too fresh (age=%v)  mark updating", it.TrackID, age)
				_ = UpdateTrackIdStatus(it.TrackID, "updating")
			} else {
				log.Printf("track_id=%s remains updating (age=%v)", it.TrackID, age)
			}
		}
	}

	log.Printf("###################################################################")
	return nil
}

// getTrackProcessLimit reads the limit from environment variable or returns default
func getTrackProcessLimit() int {
	limitStr := os.Getenv("TRACK_PROCESS_LIMIT")
	if limitStr == "" {
		log.Println("TRACK_PROCESS_LIMIT not set, defaulting to 20")
		return 20
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("Invalid TRACK_PROCESS_LIMIT value: %v, defaulting to 20", err)
		return 20
	}
	return limit
}
