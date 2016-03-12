package trend_analysis

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

/**
File: trend_analysis.go
Description: processes and identifies trending keywords
*/

// Starts processing
func Process(numPages int, flexible bool, location string, days int) error {

	for i := 1; i < numPages+1; i++ {
		theUrl, err := buildUrl(i, flexible, location)
		if err != nil {
			log.Printf("Error building url: %v", err)
			continue
		}
		resp, err := loadPage(theUrl)
		if err != nil {
			log.Printf("Error loading page: %v", err)
			continue
		}
		log.Println("Page Count: " + strconv.Itoa(resp.PageCount))
	}
	return nil
}

func processResponse(jqr *JobQueryResponse) {
	for _, result := range jqr.Results {
		rx, err := regexp.Compile("<[^>]*>")
		if err != nil {
			return
		}
		cleanString := strings.TrimSpace(rx.ReplaceAllString(result.Contents, " "))
		words := strings.Split(cleanString, " ")
		log.Printf("Sanitized: %v", cleanString)
		log.Printf("Words: %v", words)
	}
}
