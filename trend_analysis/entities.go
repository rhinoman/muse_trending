package trend_analysis

/**
File: entities.go
Description: Types for the Analyzer
*/

import (
	"sync"
	"time"
)

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

// Need a thread-safe map for storing Term Frequencies
type TFSet struct {
	TermFreq map[string]float64
	sync.RWMutex
}

//And another for storing number of docs containing a term
type DFSet struct {
	DocFreq map[string]int
	NumDocs int
	sync.RWMutex
}

type TFIDF map[string]float64

type TrendingTerm struct {
	Term  string
	Score float64
}

type TrendingTermList []TrendingTerm

//Need to override some functions to enable sorting for the TrendingTermList
func (ttl TrendingTermList) Len() int           { return len(ttl) }
func (ttl TrendingTermList) Less(i, j int) bool { return ttl[i].Score < ttl[j].Score }
func (ttl TrendingTermList) Swap(i, j int)      { ttl[i], ttl[j] = ttl[j], ttl[i] }

type StopWords map[string]bool
