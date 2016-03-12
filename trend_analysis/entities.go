package trend_analysis

/**
File: entities.go
Description: Types for the Analyzer
*/

import "time"

// Fortunately, we only have to define fields we care about here
type Job struct {
	Id              int           `json:"id"`
	Contents        string        `json:"contents"`
	Locations       []Location    `json:"locations"`
	PublicationDate time.Time     `json:"publication_date"`
	Company         CompanyRecord `json:"company"`
	Tags            []Tag         `json:"tags"`
}

type Location struct {
	Name string `json:"name"`
}

type Tag struct {
	ShortName string `json:"short_name"`
	Name      string `json:"name"`
}

type CompanyRecord struct {
	Name      string `json:"name"`
	Id        int    `json:"id"`
	ShortName string `json:"short_name"`
}

// The response object
type JobQueryResponse struct {
	Results   []Job `json:"results"`
	PageCount int   `json:"page_count"`
	PageNum   int   `json:"page"`
}
