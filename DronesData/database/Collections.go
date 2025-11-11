package database

import "go.mongodb.org/mongo-driver/mongo"

func SensorsDataCollection() *mongo.Collection {

	return GetClient().Database("droniq_flights_go").Collection("sensors_data")

}

func AllFlightsDataCollection() *mongo.Collection {
	return GetClient().Database("droniq_flights_go").Collection("flights_data")
}

func FlightLocationCollection() *mongo.Collection {
	return GetClient().Database("droniq_flights_go").Collection("flight_location")
}

func TrackIdsCollection() *mongo.Collection {
	return GetClient().Database("droniq_flights_go").Collection("track_ids")
}
