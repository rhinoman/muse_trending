package trend_analysis

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

/**
File: stop_words.go
Description: Loads a list of stop words into a map
*/

// Load our stop words
func loadStopWords(filename string) StopWords {
	if filename == "" { //No filename provided, do the default
		//Get the directory where the main executable lives
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		filename, err = filepath.Abs(dir + "/trend_analysis/stop_words.txt")
		if err != nil {
			panic(err)
		}
	}
	//open the file
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close() // Remember to close the file
	scanner := bufio.NewScanner(file)
	sw := StopWords{}
	//Read lines from the file (one stop word per line)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		sw[line] = true
	}
	return sw
}
