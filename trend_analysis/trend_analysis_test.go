package trend_analysis

/**
File: trend_analysis_test.go
Description: Unit tests for the trend analyzer
*/

import (
	"strings"
	"testing"
	"time"
)

// Tests that URL strings are sane
func TestBuildUrl(t *testing.T) {
	Init()
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
	Init()
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
	Init()
	words := []string{"f,oo.ba\nr?!", "donut", "app!le\n", "\tbana,na"}
	// The expected result, to test against
	expected := [4]string{"foobar", "donut", "apple", "banana"}
	processWords(words, false)
	controlSet.RLock()
	for _, word := range expected {
		if controlSet.DocFreq[word] != 1 {
			t.Errorf("The word %v is not present!", word)
		}
	}
	controlSet.RUnlock()
}

//Test process an API response
func TestProcessResponse(t *testing.T) {
	Init()
	jqr := GenerateJobQueryResponse()
	processResponse(jqr, 30, "")
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
	processResponse(jqr, 60, "")
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
	Init()
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
