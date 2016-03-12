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

//Test process a single word
func TestProcessWord(t *testing.T) {
	Init()
	word := "f,oo.ba\nr!?"
	processWord(word, false)
	controlSet.Mutex.Lock()
	if controlSet.Words["foobar"] == 0 {
		t.Errorf("foobar not present!")
	}
	controlSet.Mutex.Unlock()
}

//Test process list of words
func TestProcessWords(t *testing.T) {
	Init()
	words := []string{"f,oo.ba\nr?!", "donut", "app!le\n", "\tbana,na"}
	// The expected result, to test against
	expected := [4]string{"foobar", "donut", "apple", "banana"}
	processWords(words, false)
	controlSet.Mutex.Lock()
	for _, word := range expected {
		if controlSet.Words[word] != 1 {
			t.Errorf("The word %v is not present!", word)
		}
	}
	controlSet.Mutex.Unlock()
}

//Test process an API response
func TestProcessResponse(t *testing.T) {
	Init()
	jqr := GenerateJobQueryResponse()
	processResponse(jqr, 30, "")
	controlSet.Mutex.Lock()
	targetSet.Mutex.Lock()
	if controlSet.Words["bar"] != 1 {
		t.Error("bar not loaded in control set")
	}
	if targetSet.Words["bar"] != 0 {
		t.Error("bar should not be in target set")
	}
	controlSet.Mutex.Unlock()
	targetSet.Mutex.Unlock()
	processResponse(jqr, 60, "")
	targetSet.Mutex.Lock()
	if targetSet.Words["bar"] != 1 {
		t.Error("bar not loaded into target set!")
	}
	targetSet.Mutex.Unlock()
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
