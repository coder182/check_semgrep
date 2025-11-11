package model

import (
	"time"
)

type SensorsData struct {
	Serial           string    `bson:"serial"`
	Location         Location  `bson:"location"`
	Type             *string   `bson:"type,omitempty"` // Pointer for nullable field
	Name             string    `bson:"name"`
	HeightAGL        *float64  `bson:"heightAGL,omitempty"`
	HeightGeodetic   float64   `bson:"heightGeodetic"`
	LastOnline       time.Time `bson:"lastOnline"`
	Pressure         float64   `bson:"pressure"`
	Temperature      float64   `bson:"temperature"`
	IPAddress        *string   `bson:"ipAddress,omitempty"`
	WarningRadius    *float64  `bson:"warningRadius,omitempty"`
	AlertRadius      *float64  `bson:"alertRadius,omitempty"`
	Online           bool      `bson:"online"`
	State            int       `bson:"state"`
	NoiseFloorValue  float64   `bson:"noiseFloorValue"`
	NoiseFloorStatus string    `bson:"noiseFloorStatus"`
}

type Location struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"` // [longitude, latitude]
}
