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
	stopWordsFile := flag.String("stopWords", "", "list of stop words for filtering")
	days := flag.Int("days", 30, "analyze jobs for the last x days")
	numResults := flag.Int("numResults", 10, "number of trending terms to display")
	flag.Parse()
	// Initialize our analyzer
	trend_analysis.Init(*stopWordsFile)
	// Start processing data
	termList, errs := trend_analysis.Process(*numPages, *flexible, *location, *days)
	fmt.Println("Finished Processing")
	fmt.Printf("\n**** %v Errors occurred during processing ****\n", len(errs))
	//Go through the results
	fmt.Println("==== Trending Terms ====")
	if *location != "" {
		fmt.Printf("\nLocation: %v\n", *location)
	}
	fmt.Printf("Displaying the top %v trending terms\n", *numResults)
	//Use a tabwriter for pretty output
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)
	displayed := 0
	for i, term := range termList {
		if displayed > *numResults {
			break
		}
		fmt.Fprintf(w, "%v)\tTerm: %v\t\tScore: %.4f\n", i, term.Term, term.Score)
		displayed++
	}
	w.Flush()
}
