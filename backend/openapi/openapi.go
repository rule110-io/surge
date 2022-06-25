package openapi

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetAvgFee() string {
	url := "https://openapi.nkn.org/api/v1/statistics/avgtxfee"
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	result := string(body[:])
	return result
}
