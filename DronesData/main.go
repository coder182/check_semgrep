package main

import (
	"DronesData/config"
	"DronesData/controller"
	"DronesData/schedulers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	// *************************
	// Calling configuration file to load properties file

	config.LoadLocal()
	log.Printf("Properties file loaded")

	// *************************

	router := mux.NewRouter()
	controller.DronesRouter(router)

	//***************************
	// calling schedulers
	schedulers.StartAllSchedulers()

	//**********************************
	// Server ports from env file
	port := os.Getenv("SERVER_PORT")

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
