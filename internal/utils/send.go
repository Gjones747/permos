package utils

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

/*
	This file will just send the notifaction to your phone with the dates that now have open spots
	and the number of spots that are now open on that day
*/

// datesOpened should include the dates that have been opened and the new total number of spots on
// that day
func Send(userConfig *Config, datesOpened map[time.Time]int) error {
	dates := make([]time.Time, len(datesOpened))
	for date := range datesOpened {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	reqString := ""
	for _, date := range dates {
		reqString += date.Format(time.RFC822)
		reqString += ":   "
		reqString += fmt.Sprintf("%d", datesOpened[date])
		reqString += "\n"
	}

	url := "https://ntfy.sh/" + userConfig.SecretTopic

	req, err := http.NewRequest("POST", url, strings.NewReader(reqString))
	if err != nil {
		return err
	}

	req.Header.Set("Priority", "5")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return errors.New(resp.Status)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}

func SendStartMsg(userConfig *Config) error {
	url := "https://ntfy.sh/" + userConfig.SecretTopic

	req, err := http.NewRequest("POST", url, strings.NewReader("Watching for new permits!"))
	if err != nil {
		return err
	}

	req.Header.Set("Priority", "3")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return errors.New(resp.Status)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

func SendErr(userConfig *Config, err error) error {
	url := "https://ntfy.sh/" + userConfig.SecretTopic

	req, err := http.NewRequest("POST", url, strings.NewReader(err.Error()))
	if err != nil {
		return err
	}

	req.Header.Set("Priority", "1")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return errors.New(resp.Status)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil

}
