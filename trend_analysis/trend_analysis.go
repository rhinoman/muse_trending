package trend_analysis

import (
	"log"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"
)

/**
File: trend_analysis.go
Description: processes and identifies trending keywords
*/

var targetSet TFSet
var controlSet DFSet
var tagRegex *regexp.Regexp
var punctRegex *regexp.Regexp

// Initializes some variables
func Init() {
	// pre-compile some regular expressions
	tagRegex = regexp.MustCompile("<[^>]*>")
	punctRegex = regexp.MustCompile("[.,?!;:()\n\r\t]")
	// Initialize our data structures
	targetSet = TFSet{
		TermFreq: make(map[string]float64),
	}
	controlSet = DFSet{
		DocFreq: make(map[string]int),
	}

}

// Starts processing
func Process(numPages int, flexible bool, location string, days int) (TrendingTermList, []error) {
	// Keep track of our errors
	errorList := []error{}
	// Start pulling data
	for i := 0; i < numPages; i++ {
		theUrl, err := buildUrl(i, flexible, location)
		if err != nil {
			log.Printf("Error building url: %v", err)
			errorList = append(errorList, err)
			continue
		}
		resp, err := loadPage(theUrl)
		if err != nil {
			log.Printf("Error loading page: %v", err)
			errorList = append(errorList, err)
			continue
		}
		processResponse(resp, days, location)
	}
	tfidf := computeTFIDF()

	return sortResults(tfidf), errorList
}

// Sorts the final results into reverse order
func sortResults(results TFIDF) TrendingTermList {
	ttl := make(TrendingTermList, len(results))
	i := 0
	for k, v := range results {
		ttl[i] = TrendingTerm{k, v}
		i++
	}
	sort.Sort(sort.Reverse(ttl))
	return ttl
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
	if target {
		occurrences := make(map[string]int)
		for _, word := range words {
			if cleanedWord := string(punctRegex.ReplaceAllString(word, "")); cleanedWord != "" {
				occurrences[cleanedWord] += 1
			}
		}
		//Compute term frequencies for each word
		targetSet.Lock()
		docSize := float64(len(words))
		for k, v := range occurrences {
			targetSet.TermFreq[k] = float64(v) / docSize
		}
		targetSet.Unlock()
	} else {
		//Simpler, here we're only interested in the number of documents which contain a word
		hasOccurred := make(map[string]bool)
		controlSet.Lock()
		//Increase our document count
		controlSet.NumDocs += 1
		for _, word := range words {
			if cleanedWord := string(punctRegex.ReplaceAllString(word, "")); cleanedWord != "" && !hasOccurred[cleanedWord] {
				hasOccurred[cleanedWord] = true
				controlSet.DocFreq[cleanedWord] += 1
			}
		}
		controlSet.Unlock()
	}
}

// Compute the TF-IDF of our target set
func computeTFIDF() TFIDF {
	//Create a map to hold the TF-IDF for our keywords
	tfidf := TFIDF{}

	//Now let's go through our results
	targetSet.RLock()
	controlSet.RLock()
	//compute the tf-idf for each word in the target set
	for word, targetNum := range targetSet.TermFreq {
		df := float64(controlSet.DocFreq[word])
		n := float64(controlSet.NumDocs)
		idf := math.Log(n / (df + 1.0))
		tfidf[word] = float64(targetNum) * idf
	}
	controlSet.RUnlock()
	targetSet.RUnlock()
	return tfidf
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
