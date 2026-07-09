package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
	This file just gets the user configuration for how they want to check permits
*/

var CONFIG_FILE_NAME string = "config.json"

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

func resolveConfigLocation() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	resolved, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", err
	}

	return filepath.Dir(resolved), nil
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

	exeDir, err := resolveConfigLocation()
	if err != nil {
		return Config{}, err
	}

	resolvedPath := filepath.Join(exeDir, CONFIG_FILE_NAME)

	data, err := os.ReadFile(resolvedPath)
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
