package main

import (
	"fmt"
	"log"
	"maps"
	"permos/internal/utils"
	"time"
)

var DEFAULT_CONFIG utils.Config = utils.Config{
	Url:         "https://www.recreation.gov/api/permits/4675309/divisions/999/availability?start_date=2026-07-04T07:00:00.000Z&end_date=2026-08-31T00:00:00.000Z&commercial_acct=false&is_lottery=false",
	RequiredNum: 0,
	Dates:       nil,
	SecretTopic: "permos-noties",
	Timeout:     100000,
}

func main() {
	log.Println("Hello Welcome to permos")

	userConfig, err := utils.GetConfig()
	if err != nil {
		log.Printf("Error processing config: %s - Using defaults", err)
		userConfig = DEFAULT_CONFIG
	}

	availabilityMap := make(map[time.Time]int)
	prevAvailabilityMap := make(map[time.Time]int)

	err = utils.Fetch(userConfig.Url, availabilityMap)
	if err != nil {
		log.Panicf("Failed to fetch address in config, please include a valid endpoint")
	}

	maps.Copy(availabilityMap, prevAvailabilityMap)

	fmt.Println("Current Openings:")
	for key, val := range availabilityMap {
		fmt.Printf("%s: %d\n", key, val)
	}

	fmt.Println("")
	log.Println("Starting Notify loop...")

	utils.SendStartMsg(&userConfig)
	log.Println("Watching for permits...")

	for {

		time.Sleep(time.Duration(userConfig.Timeout) * time.Minute)

		err = utils.Fetch(userConfig.Url, availabilityMap)
		if err != nil {
			log.Panicf("Failed to fetch address in config, please include a valid endpoint")
		}

		diffMap := make(map[time.Time]int)

		for key, val := range availabilityMap {
			if val > prevAvailabilityMap[key] {
				diffMap[key] = val
			}
		}

		if len(diffMap) != 0 {
			utils.Send(&userConfig, diffMap)
		}
	}
}
