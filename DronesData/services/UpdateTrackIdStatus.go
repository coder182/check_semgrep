package services

import (
	"DronesData/database"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// UpdateTrackStatus updates the status of a track ID in the track_ids collection.
func UpdateTrackIdStatus(trackID string, status string) error {
	collection := database.TrackIdsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"track_id": trackID}
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Failed to update status for track_id %s: %v", trackID, err)
		return err
	}

	if result.MatchedCount == 0 {
		log.Printf("No document matched for track_id %s", trackID)
	} else {
		log.Printf("Track id status updated successful")
	}
	return nil
}
