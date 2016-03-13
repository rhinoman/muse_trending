package trend_analysis

/**
File: trend_analysis_test.go
Description: Unit tests for the trend analyzer
*/

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Set up prior to a test
func beforeTest() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	//The location of the stop words file is not the same when running in test
	//So need to figure it out here
	swFile, err := filepath.Abs(dir + "/../stop_words.txt")
	if err != nil {
		panic(err)
	}
	Init(swFile)
}

// Tests that URL strings are sane
func TestBuildUrl(t *testing.T) {
	beforeTest()
	theUrl, err := buildUrl(42, true, "The Moon")
	if err != nil {
		t.Error(err)
	}
	t.Logf("The URL: %v", theUrl)
	if !strings.Contains(theUrl, apiLocation) {
		t.Error("apiLocation not in URL!")
	}
	if !strings.Contains(theUrl, "page=42") {
		t.Errorf("Page number should be 42!")
	}
	if !strings.Contains(theUrl, "flexible=true") {
		t.Errorf("Flexible should be true!")
	}
	if !strings.Contains(theUrl, "The+Moon") {
		t.Error("Location is wrong!")
	}

}

// Requires internet connection
func TestLoadPage(t *testing.T) {
	//Skipping so we don't hit the network
	//Comment out this line to run test
	t.SkipNow()
	beforeTest()
	jqr, err := loadPage(apiLocation + "?page=1")
	if err != nil {
		t.Error(err)
	}
	if len(jqr.Results) <= 0 {
		t.Error("Results are empty!")
	}
	if jqr.PageCount <= 0 {
		t.Error("Page Count is 0!")
	}
	if jqr.PageNum != 1 {
		t.Error("Page should be 1!")
	}
}

//Test process list of words
func TestProcessWords(t *testing.T) {
	beforeTest()
	words := []string{"f,oo.ba\nr?!", "donut ", "ap p!le\n", "\tbana,na", "and"}
	// The expected result, to test against
	expected := [4]string{"foobar", "donut", "apple", "banana"}
	processWords(words, false)
	controlSet.RLock()
	for _, word := range expected {
		if controlSet.DocFreq[word] != 1 {
			t.Errorf("The word %v is not present!", word)
		}
	}
	//Make sure the stop word filter is doing something
	if controlSet.DocFreq["and"] != 0 {
		t.Error("and is present and it shouldn't be")
	}
	controlSet.RUnlock()
}

//Test process an API response
func TestProcessResponse(t *testing.T) {
	beforeTest()
	jqr := GenerateJobQueryResponse()
	wg.Add(1)
	go processResponse(jqr, 30, "")
	wg.Wait()
	controlSet.RLock()
	targetSet.RLock()
	if controlSet.DocFreq["bar"] != 1 {
		t.Error("bar not loaded in control set")
	}
	if targetSet.TermFreq["bar"] != 0 {
		t.Error("bar should not be in target set")
	}
	controlSet.RUnlock()
	targetSet.RUnlock()
	wg.Add(1)
	go processResponse(jqr, 60, "")
	wg.Wait()
	targetSet.RLock()
	if targetSet.TermFreq["bar"] <= 0 {
		t.Logf("TermFreq: %v", targetSet.TermFreq["bar"])
		t.Error("bar not loaded into target set!")
	}
	targetSet.RUnlock()
}

//Test TFIDF
func TestComputeTFIDF(t *testing.T) {
	//Generate our sets
	beforeTest()
	controlSet.Lock()
	controlSet.DocFreq["foo"] = 5
	controlSet.DocFreq["bar"] = 2
	controlSet.DocFreq["rhinoceros"] = 1
	controlSet.NumDocs = 5
	controlSet.Unlock()
	targetSet.Lock()
	targetSet.TermFreq["foo"] = 4
	targetSet.TermFreq["bar"] = 2
	targetSet.TermFreq["rhinoceros"] = 3
	targetSet.NumTerms = 9
	targetSet.Unlock()
	tfidf := computeTFIDF()
	t.Logf("\ntfidf: %v\n", tfidf)
	if !(tfidf["foo"] < tfidf["bar"]) && (tfidf["bar"] < tfidf["rhinoceros"]) {
		t.Error("TFIDF calculation is wrong!")
	}

}

//Test target set determination function
func TestInTarget(t *testing.T) {
	job := GenerateJob(42, "foo")
	if inTarget(job, 10, "") {
		t.Error("Should not be in target")
	}
	if !inTarget(job, 60, "") {
		t.Error("Should be in target!")
	}
	if inTarget(job, 10, "Moon") {
		t.Error("Should not be in target!")
	}
	if inTarget(job, 60, "Mars") {
		t.Error("Should not be in target!")
	}
	if !inTarget(job, 60, "Moon") {
		t.Error("Should be in target!")
	}
}

// --- Helper Functions ---

func GenerateJobQueryResponse() *JobQueryResponse {
	var jobList = make([]Job, 10)
	for i := 0; i < 10; i++ {
		content := "<h1>Lorem ipsum</h1>\n"
		if i%2 == 0 {
			content += "<p>foo</p>"
		}
		if i == 5 {
			content += "<h2>bar</h2>"
		}
		job := GenerateJob(i, content)
		jobList[i] = job
	}
	return &JobQueryResponse{
		Results:   jobList,
		PageCount: 1,
		PageNum:   1,
	}
}

func GenerateJob(id int, contents string) Job {
	now := time.Now().UTC()
	return Job{
		Id:              id,
		Contents:        contents,
		PublicationDate: now.AddDate(0, 0, -40),
		Locations:       []Location{Location{Name: "Moon"}},
	}

}
