package main

/**
File: main.go
Description: Main entry point for the muse_trending application
*/

import (
	"flag"
	"fmt"
	"github.com/rhinoman/muse_trending/trend_analysis"
	"os"
	"text/tabwriter"
)

func main() {
	//Get the command line arguments
	numPages := flag.Int("numPages", 10, "number of pages to process")
	flexible := flag.Bool("flexible", true, "include jobs with flexible location in analysis")
	location := flag.String("location", "", "analyze jobs for a specific location")
	days := flag.Int("days", 30, "analyze jobs for the last x days")
	numResults := flag.Int("numResults", 10, "number of ternding terms to display")
	flag.Parse()
	// Initialize our analyzer
	trend_analysis.Init()
	// Start processing data
	termList, errs := trend_analysis.Process(*numPages, *flexible, *location, *days)
	fmt.Println("Finished Processing")
	fmt.Printf("\n**** %v Errors occurred during processing ****\n", len(errs))
	//Go through the results
	fmt.Println("==== Trending Terms ====")
	fmt.Printf("Displaying the top %v trending terms\n", *numResults)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	displayed := 0
	for i, term := range termList {
		if displayed > *numResults {
			break
		}
		fmt.Fprintf(w, "%v:\tTerm: %v\tScore: %.4f\n", i, term.Term, term.Score)
		displayed++
	}
	w.Flush()
}
