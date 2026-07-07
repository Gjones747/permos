package main

import (
	"log"
	"permos/internal/api"
	"time"
)

var HELENS_ENDPOINT string = "https://www.recreation.gov/api/permits/4675309/divisions/999/availability?start_date=2026-07-04T07:00:00.000Z&end_date=2026-08-31T00:00:00.000Z&commercial_acct=false&is_lottery=false"

func main() {
	log.Println("Hello Welcome to permos")

	availabilityMap := make(map[time.Time]int)

	api.Fetch(HELENS_ENDPOINT, availabilityMap)
}
