package services

import (
	"DronesData/database"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetLastTrackingIDs(isFristRun bool) ([]string, error) {
	log.Printf("Reading last track id from track_id collection")
	trackIdCollection := database.TrackIdsCollection()
	allFlightcollection := database.AllFlightsDataCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize filter - always require trackid to exist
	filter := bson.M{"track_id": bson.M{"$exists": true, "$ne": ""}}

	if !isFristRun {
		// Step 1: Get last inserted _id from TrackIdsCollection
		opts := options.FindOne().SetSort(bson.D{{"_id", -1}}).SetProjection(bson.M{"_id": 1})
		var latestDoc struct {
			ID interface{} `bson:"_id"`
		}

		err := trackIdCollection.FindOne(ctx, bson.M{}, opts).Decode(&latestDoc)
		if err == nil && latestDoc.ID != nil {
			// Step 2: Filter documents with _id > latestDoc.ID
			filter["_id"] = bson.M{"$gt": latestDoc.ID}
			log.Printf("Filtering AllFlightsDataCollection for records with _id > %v", latestDoc.ID)
		} else {
			if err != nil && err != mongo.ErrNoDocuments {
				log.Printf("Error finding latest _id: %v", err)
				return nil, fmt.Errorf("failed to get latest _id: %w", err)
			}
			log.Println("No existing records found in trackIdCollection - fetching all records")
		}
	}

	findOpts := options.Find().SetProjection(bson.M{
		"track_id": 1,
		"end_time": 1,
		"_id":      1,
	})

	// Step 3: Fetch filtered records from AllFlightsDataCollection
	cursor, err := allFlightcollection.Find(ctx, filter, findOpts)
	if err != nil {
		log.Printf("Failed to fetch from All Flights Collection: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// end_time is an embedded doc: { time: <Date> }
	type resultDoc struct {
		TrackID string `bson:"track_id"`
		EndTime struct {
			Time time.Time `bson:"time"`
		} `bson:"end_time"`
	}

	// var results []struct {
	// 	TrackID string `bson:"track_id"`
	// 	EndTime struct {
	//         Time time.Time `bson:"time"`
	//     } `bson:"end_time"`
	// }

	var results []resultDoc

	if err = cursor.All(ctx, &results); err != nil {
		log.Printf("Cursor iteration failed: %v", err)
		return nil, err
	}

	log.Printf("Found %d new track IDs", len(results))

	// List of successful IDs with end_time
	var trackData []struct {
		TrackID string
		EndTime time.Time
	}

	for _, item := range results {
		if item.TrackID == "" {
			log.Println("Skipping empty track ID")
			continue
		}
		trackData = append(trackData, struct {
			TrackID string
			EndTime time.Time
		}{
			TrackID: item.TrackID,
			EndTime: item.EndTime.Time,
		})
	}

	// need correct name as trackidsService
	err = InsertTrackIDsToCollection(trackData)
	if err != nil {
		log.Printf("Failed to insert track IDs into track_ids collection: %v", err)
	}
	log.Printf("#####################################################")
	return nil, nil
}
