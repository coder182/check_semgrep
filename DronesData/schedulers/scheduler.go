package schedulers

import (
	"DronesData/services"
	"log"
	"os"
	"sync"

	"github.com/robfig/cron/v3"
)

// isFirstRun had problem if you run the application and then stop,
// the flag will reset to false and it will duplicate the complete data in TrackId collection,
// to handle this we are using docker volume path
var (
	stateInit     sync.Once
	isFirstRun    = true
	stateFilePath = "/data/schedular_first_run"
)

func StartAllSchedulers() {
	c := cron.New()

	// FetchAndStoreDrones every 10 min
	c.AddFunc("*/3 * * * *", func() {
		log.Println("Sensors API Scheduler initialized successfully")
		if err := services.FetchAndStoreSensors(); err != nil {
			log.Printf("RetryFailedTrackIds failed: %v", err)
		}
		log.Println("##################################################")
	})

	//  StoreAllFlightsData every 10 min
	c.AddFunc("*/3 * * * *", func() {
		log.Println("Pagination API Scheduler initialized successfully...")
		if err := services.StoreAllFlightsData(); err != nil {
			log.Printf("FetchNewTrackIds failed: %v", err)
		}
		log.Println("##################################################")
	})

	// GetLastTrackingIDs every 10 min
	c.AddFunc("*/3 * * * *", func() {
		stateInit.Do(initSchedularState)
		log.Println("Getting last Track Id  Scheduler initialized successfully...")
		if _, err := services.GetLastTrackingIDs(isFirstRun); err != nil {
			log.Printf("CleanFailedTrackIds failed: %v", err)
		} else if isFirstRun {
			isFirstRun = false
			persisSchedularState()
		}
		log.Println("##################################################")
	})

	//  ProcessTrackIds every 3 min
	c.AddFunc("*/3 * * * *", func() {
		log.Println("Processing Track Id  Scheduler initialized successfully...")
		if err := services.ProcessTrackIds(); err != nil {
			log.Printf(" ProcessTrackIds failed: %v", err)
		}

		log.Println("##################################################")
	})

	c.Start()
	log.Println(" All schedulers started...")

	// Prevent app from exiting
	select {}
}

func initSchedularState() {
	if _, err := os.Create(stateFilePath); err == nil {
		isFirstRun = false
		log.Printf("Schedular: Resuming from previous state")
	}
}

func persisSchedularState() {
	if file, err := os.Create(stateFilePath); err == nil {
		file.Close()
		log.Printf("Schedular: persisted first-run state")
	} else {
		log.Printf("Schedular: Warning - failed to persist state: %v", err)
	}
}
