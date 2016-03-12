package trend_analysis

/**
File: trend_analysis_test.go
Description: Unit tests for the trend analyzer
*/

import (
	"strings"
	"testing"
)

// Tests that URL strings are sane
func TestBuildUrl(t *testing.T) {
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

func TestProcessResponse(t *testing.T) {
	jqr := GenerateJobQueryResponse()
	processResponse(jqr)
}

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

	return Job{
		Id:       id,
		Contents: contents,
	}

}
