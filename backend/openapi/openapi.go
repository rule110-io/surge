package openapi

import (
	"io/ioutil"
	"net/http"
	"time"
)

func GetAvgFee() (string, error) {
	url := "https://openapi.nkn.org/api/v1/statistics/avgtxfee"
	client := http.Client{
		Timeout: time.Second * 10, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		return "", getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", readErr
	}
	result := string(body[:])
	return result, nil
}
