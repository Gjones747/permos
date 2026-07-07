package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

/*
	Purpose of this file is to just go out and fetch the park service permit data sources and put
	that data into a readable map
*/

type Response struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	PermitId          string             `json:"permit_id"`
	NextAvailableDate string             `json:"next_available_date"`
	DateAvailability  map[string]Day     `json:"date_availability"`
	QuotaTypeMaps     QuotaUsageByMember `json:"quota_type_maps"`
}

type QuotaUsageByMember map[string]Day

type Day struct {
	Total         int  `json:"total"`
	Remaining     int  `json:"remaining"`
	ShowWalkUp    bool `json:"show_walkup"`
	IsSecretQouta bool `json:"is_secret_qouta"`
}

// This method takes an endpoint to a national parks permit resource and will parse the response
// into a map from date -> availability
func Fetch(endPoint string, availableMap map[time.Time]int) error {

	responseJSON := Response{}

	req, err := http.NewRequest("GET", endPoint, nil)

	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:150.0) Gecko/20100101 Firefox/150.0")

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		log.Panic(err)
		return err
	}

	responseBytes, err := io.ReadAll(response.Body)

	err = json.Unmarshal(responseBytes, &responseJSON)

	if err != nil {
		log.Panic(err)
		return err
	}

	fmt.Printf("%s\n", responseJSON.Payload.NextAvailableDate)

	for key, val := range responseJSON.Payload.DateAvailability {

		permitTime, err := time.Parse(time.RFC3339, key)

		if err != nil {
			log.Panic(err)
		}

		fmt.Println(key + " " + strconv.Itoa(val.Remaining))

		availableMap[permitTime] = val.Remaining
	}

	return nil
}
