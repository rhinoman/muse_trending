package trend_analysis

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

/**
File: net.go
Description: Contains network request code
*/

// Location of The Muse Jobs API
const apiLocation = "https://api-v2.themuse.com/jobs"

// Our HTTP(s) client
var client = http.Client{Timeout: 5 * time.Second}

// Sends request to The Muse API
func loadPage(queryString string) (*JobQueryResponse, error) {
	req, err := http.NewRequest("GET", queryString, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Make sure we close the response body when we're done
	defer resp.Body.Close()
	jqr := JobQueryResponse{}
	err = parseBody(resp, &jqr)
	return &jqr, err
}

// Creates a URL for querying a page
func buildUrl(pageNum int, flexible bool, location string) (string, error) {
	museUrl, err := url.Parse(apiLocation)
	if err != nil {
		return "", err
	}
	//Set the query parameters
	params := url.Values{}
	params.Set("page", strconv.Itoa(pageNum))
	params.Add("flexible", strconv.FormatBool(flexible))
	if location != "" {
		params.Add("location", location)
	}
	museUrl.RawQuery = params.Encode()
	return museUrl.String(), nil
}

// unmarshalls a JSON response body
func parseBody(resp *http.Response, o interface{}) error {
	err := json.NewDecoder(resp.Body).Decode(&o)
	if err != nil {
		resp.Body.Close()
		return err
	} else {
		return resp.Body.Close()
	}
}
