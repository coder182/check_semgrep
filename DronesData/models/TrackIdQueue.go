package model

import "time"

type TrackIdQueue struct {
	TrackID      string     `bson:"track_id" json:"track_id"`
	Status       string     `bson:"status" json:"status"` // pending, processing, done, failed
	InsertedDate CustomTime `bson:"inserted_date" json:"inserted_date"`
	End_time     time.Time  `bson:"end_time" json:"end_time"`
}
