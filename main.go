package main

import (
	"fmt"
	. "github.com/ChuckHa/Cuttlefish/cuttlefish"
	"regexp"
)

func main() {
	regex := regexp.MustCompile("<a.*?href=[\"'](http.*?)[\"']")

	curl := make(chan []byte)
	csite := make(chan Site)
	done := make(chan struct{})

	// Give our crawler a place to start.
	go Seed(curl)

	// Keeps track of which urls we have visted.
	visited := make(map[string]int)

	// Start the throttled crawling.
	go ThrottledCrawl(curl, csite, done, visited)

	// Main loop that never exits and blocks on the data of a page.
	for {
		fmt.Println("blocking on getting a site")
		site := <-csite
		fmt.Println("Got a site")
		go GetUrls(curl, site, regex)
	}
}
