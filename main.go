package main

/**
File: main.go
Description: Main entry point for the muse_trending application
*/

import (
	"flag"
	"github.com/rhinoman/muse_trending/trend_analysis"
)

func main() {
	//Get the command line arguments
	numPages := flag.Int("numPages", 10, "number of pages to process")
	flexible := flag.Bool("flexible", true, "include jobs with flexible location in analysis")
	location := flag.String("location", "", "analyze jobs for a specific location")
	days := flag.Int("days", 90, "analyze jobs for the last x days")
	flag.Parse()
	// Initialize our analyzer
	trend_analysis.Init()
	// Start processing data
	trend_analysis.Process(*numPages, *flexible, *location, *days)
}
