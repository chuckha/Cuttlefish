package cuttlefish

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Site struct {
	Url, Body []byte
}

// GetUrl will make an HTTP GET request, build a site object and put it on a channel.
// It will send a message on the stop channel after the function finishes.
func GetUrl(url []byte, csite chan Site, done chan struct{}) {
	resource := string(url)
	defer func () {
		done <- struct{}{}
	}()
	resp, err := http.Get(resource)
	if err != nil {
		fmt.Println("We have an error!: ", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Getting %v\n", resource)
	body, _ := ioutil.ReadAll(resp.Body)
	site := Site{url, body}
	csite <- site
}

// ThrottledCrawl will limit the number of goroutines making requests.
// It will listen to a URL channel and spawn a goroutine for each URL.
// It manages the number of goroutines using a stop channel.
// This function does not return and should be used as a goroutine.
func ThrottledCrawl(curl chan []byte, csite chan Site, done chan struct{}, visited map[string]int) {
	maxGos := 10
	numGos := 0
	for {
		if numGos > maxGos {
			<-done
			numGos -= 1
		}
		url := string(<-curl)
		if _, ok := visited[url]; !ok {
			go GetUrl([]byte(url), csite, done)
			numGos += 1
		}
		visited[url] += 1
	}
}

// Seed starts the crawling process by feeding the URL channel a URL.
func Seed(curl chan []byte) {
	curl <- []byte("https://news.ycombinator.com")
}

// GetUrls parses a site object and looks for links to sites.
func GetUrls(curl chan []byte, site Site, regex *regexp.Regexp) {
	matches := regex.FindAllSubmatch(site.Body, -1)
	for _, match := range matches {
		curl <- match[1]
	}
}

