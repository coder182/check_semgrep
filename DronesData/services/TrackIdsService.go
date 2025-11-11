package services

import (
	"DronesData/database"
	model "DronesData/models"
	"context"
	"log"
	"time"
)

// InsertTrackIDsToCollection takes a list of track IDs and inserts them into the track_ids collection with status and inserted_date.
func InsertTrackIDsToCollection(trackData []struct {
	TrackID string
	EndTime time.Time
}) error {
	collection := database.TrackIdsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var docs []interface{}
	currentTime := time.Now()

	for _, data := range trackData {
		doc := model.TrackIdQueue{
			TrackID:      data.TrackID,
			Status:       "pending",
			InsertedDate: model.CustomTime{Time: currentTime},
			End_time:     data.EndTime,
		}
		docs = append(docs, doc)
	}

	if len(docs) == 0 {
		log.Println("No track IDs to insert")
		return nil
	}

	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Printf("Failed to insert track IDs into track_ids collection: %v", err)
		return err
	}

	log.Printf("Inserted %d track IDs into track_ids collection", len(docs))
	return nil
}
