package model

type FlightLocation struct {
	TrackID             string        `json:"track_id" bson:"track_id"`
	VehicleLocation     GeoPoint      `json:"vehicle_location" bson:"vehicle_location"`
	VehicleHome         GeoPoint      `json:"vehicle_home" bson:"vehicle_home"`
	PilotLocation       GeoPoint      `json:"pilot_location" bson:"pilot_location"`
	Trace               LineString    `json:"trace" bson:"trace"`
	VehicleSerialNumber string        `json:"vehicle_serial_number" bson:"vehicle_serial_number"`
	HeightAGL           float64       `json:"height_agl" bson:"height_agl"`
	HeightAGLHistory    HeightHistory `json:"height_agl_history" bson:"height_agl_history"`
	TrackDirection      float64       `json:"track_direction" bson:"track_direction"`
	GroundSpeedDir      *float64      `json:"ground_speed_dir" bson:"ground_speed_dir"`
	GroundSpeedSize     *float64      `json:"ground_speed_size" bson:"ground_speed_size"`
	VehicleModel        int           `json:"vehicle_model" bson:"vehicle_model"`
	VehicleModelName    string        `json:"vehicle_model_name" bson:"vehicle_model_name"`
	SourceType          int           `json:"source_type" bson:"source_type"`
	OperatorID          *string       `json:"operator_id" bson:"operator_id"`
	VehicleMake         *string       `json:"vehicle_make" bson:"vehicle_make"`
	SourceSerials       []string      `json:"source_serials" bson:"source_serials"`
	FriendStatus        int           `json:"friend_status" bson:"friend_status"`
	AlertStatus         *string       `json:"alert_status" bson:"alert_status"`
	Distance            float64       `json:"distance" bson:"distance"`
	Callsign            string        `json:"callsign" bson:"callsign"`
	ActiveDetection     *bool         `json:"active_detection" bson:"active_detection"`
}

type GeoPoint struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type LineString struct {
	Type        string      `json:"type" bson:"type"`
	Coordinates [][]float64 `json:"coordinates" bson:"coordinates"`
}

type HeightHistory struct {
	Name   string    `json:"name" bson:"name"`
	Times  []string  `json:"times" bson:"times"`
	Values []float64 `json:"values" bson:"values"`
}
