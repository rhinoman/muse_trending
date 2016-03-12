package trend_analysis

import (
	"log"
	"regexp"
	"strings"
	"sync"
	"time"
)

/**
File: trend_analysis.go
Description: processes and identifies trending keywords
*/

var targetSet WordSet
var controlSet WordSet
var tagRegex *regexp.Regexp
var punctRegex *regexp.Regexp

// Initializes some variables
func Init() {
	// pre-compile some regular expressions
	tagRegex = regexp.MustCompile("<[^>]*>")
	punctRegex = regexp.MustCompile("[.,?!;:()\n\r\t]")
	// Initialize our data structures
	targetSet = WordSet{
		Words: make(map[string]int),
		Mutex: &sync.Mutex{},
	}
	controlSet = WordSet{
		Words: make(map[string]int),
		Mutex: &sync.Mutex{},
	}

}

// Starts processing
func Process(numPages int, flexible bool, location string, days int) error {
	// Start pulling data
	for i := 0; i < numPages; i++ {
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
		processResponse(resp, days, location)
	}
	return nil
}

// Process a single API response
func processResponse(jqr *JobQueryResponse, days int, location string) {

	for _, result := range jqr.Results {
		inTarget := inTarget(result, days, location)
		// Trim leading and trailing whitespace and strip HTML tags
		cleanString := strings.TrimSpace(tagRegex.ReplaceAllString(result.Contents, " "))
		// Split resulting string into an array of individual words
		words := strings.Split(cleanString, " ")
		processWords(words, inTarget)
	}
}

// Process a list of words
func processWords(words []string, target bool) {
	for _, word := range words {
		processWord(word, target)
		// If the processed word is valid, add it to the results
	}
}

// Process a single word
func processWord(word string, target bool) {
	//First, remove any punctuation
	cleanedWord := string(punctRegex.ReplaceAllString(word, ""))
	//TODO: Then, check if this word is in our list of stopwords
	if target { //response is in our 'target' set
		targetSet.Mutex.Lock()
		targetSet.Words[cleanedWord] += 1
		targetSet.Mutex.Unlock()
	} else {
		controlSet.Mutex.Lock()
		controlSet.Words[cleanedWord] += 1
		controlSet.Mutex.Unlock()
	}
}

// Test if this job's words belong in the target set or the control set
func inTarget(job Job, days int, location string) bool {
	target := false
	now := time.Now().UTC()
	jobTime := job.PublicationDate.UTC()
	//Test if this job is in our time range
	if jobTime.AddDate(0, 0, days).After(now) {
		target = true
	}
	//If we have a location, we can further narrow our target set
	if target && (location != "") {
		target = false
		for _, loc := range job.Locations {
			if loc.Name == location {
				target = true
				break
			}
		}
	}
	return target
}
