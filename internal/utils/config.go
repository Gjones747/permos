package utils

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"
)

/*
	This file just gets the user configuration for how they want to check permits
*/

var CONFIG_FILE_PATH string = "./config.json"

type Config struct {
	Url         string       `json:"url"`
	RequiredNum int          `json:"required_num"`
	Dates       []CustomTime `json:"dates"`
	SecretTopic string       `json:"secret_topic"`
	Timeout     int          `json:"timeout"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	time_string := strings.Trim(string(b), "\"")

	if time_string == "" {
		return
	}

	ct.Time, err = time.Parse(time.DateOnly, time_string)
	return err
}

func GetConfig() (Config, error) {
	data, err := os.ReadFile(CONFIG_FILE_PATH)
	if err != nil {
		return Config{}, err
	}

	userConfig := Config{}

	log.Printf("here")

	err = json.Unmarshal(data, &userConfig)
	if err != nil {
		return Config{}, err
	}

	return userConfig, err
}
