package main

import (
	. "github.com/ChuckHa/Cuttlefish/cuttlefish"
	"regexp"
)

func main() {
	regex := regexp.MustCompile("<a.*?href=[\"'](http.*?)[\"']")

	curl := make(chan []byte)
	csite := make(chan Site)
	death := make(chan string)

	// Give our crawler a place to start.
	go Seed(curl)

	// Keeps track of which urls we have visted.
	visited := make(map[string]int)

	// Start the throttled crawling.
	go ThrottledCrawl(curl, csite, death, visited)

	// Main loop that never exits and blocks on the data of a page.
	for {
		site := <-csite
		go GetUrls(curl, site, regex)
	}
}
