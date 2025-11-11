package model

type AllFlightsResponse struct {
	Total int           `json:"total"`
	Limit *int          `json:"limit"`
	Skip  *int          `json:"skip"`
	Data  []DronesTrack `json:"data"`
}

type DronesTrack struct {
	TrackID             string     `json:"track_id" bson:"track_id"`
	VehicleSerialNumber string     `json:"vehicle_serial_number"`
	StartTime           CustomTime `json:"start_time"`
	EndTime             CustomTime `json:"end_time" bson:"end_time"`
	VehicleModel        *int       `json:"vehicle_model"`
	SourceType          int        `json:"source_type"`
	SourceTypes         []int      `json:"source_types"`
	SourceSerials       []string   `json:"source_serials"`
	MaxAltitude         float64    `json:"max_altitude"`
	FriendStatus        int        `json:"friend_status"`
	AlertStatus         string     `json:"alert_status"`
}
