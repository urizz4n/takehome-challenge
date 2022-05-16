package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"takehome-challenge/houseresponse"
	"time"
)

// Fetches and checks if the page works
func Fetch(url string, maxRetries, i int, result *houseresponse.HouseResponse) {

	pageUrl := url + "?page=" + strconv.Itoa(i)

	var resp *http.Response
	connected := false

	resp, connected = BackOff(maxRetries, pageUrl, resp)

	if connected {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading the response body")
		}

		// Parses the JSON
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatal("Can not unmarshal JSON")
		}
	}

	// If it tried too many times, continue with the program and ignore this page (we assume the site is not working)

}

func BackOff(maxRetries int, URL string, resp *http.Response) (*http.Response, bool) {

	connected := false
	retries := 1
	backoff := 0.5

	// Retries until we connect  to the page url.
	for !connected && retries <= maxRetries {

		var err error

		resp, err = http.Get(URL)
		if err != nil {
			fmt.Printf("Couldn't get %s, retrying...\n", URL)
		}

		if resp.StatusCode != 200 {
			fmt.Printf("--Response code from %s is not 200, total retries: %d...\n", URL, retries)

			// Waits so it doesn't overload the page
			time.Sleep(time.Duration(backoff) * time.Second)

			// updates the backoff and retries
			backoff = backoff * 2
			retries = retries + 1
		} else {
			connected = true
		}
	}

	if !connected {
		fmt.Printf("Too many retries (%d) when trying to fetch %s (photoUrl), ignoring the URL", retries, URL)
	}

	return resp, connected
}
