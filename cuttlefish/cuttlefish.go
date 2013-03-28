package cuttlefish

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

func RobotsUrl (resource string) string {
	// FIXME probably do something with the error
	u, _ := url.Parse(resource)
	return fmt.Sprintf("%v://%v/robots.txt", u.Scheme, u.Host)
}

type Site struct {
	Url, RobotsUrl string
	Body []byte
	done chan struct{}
}

// Constructs a new Site object
func NewSite (url string, done chan struct{}) (Site, error) {
	site := Site{Url: url, RobotsUrl: RobotsUrl(url), Body: nil,done: done}
	err := site.GetBody()
	return site, err
}

func (s Site) GetBody() (error) {
	defer func(){s.done <- struct{}{}}()
	resp, err := http.Get(s.Url)
	if err != nil {
		fmt.Println("Error on: ", string(s.Url))
		return err
	}
	defer resp.Body.Close()
	s.Body, _ = ioutil.ReadAll(resp.Body)
	return nil
}

// GetUrl will make an HTTP GET request, build a site object and put it on a channel.
// It will send a message on the stop channel after the function finishes.
func GetUrl(url string, csite chan Site, done chan struct{}) {
	site, _ := NewSite(url, done)
	fmt.Println("blocking on putting things on csite")
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
		select {
		case <-done:
			numGos -= 1
		case url := <-curl:
			surl := string(url)
			if numGos < maxGos {
				numGos += 1
				if _, ok := visited[surl]; !ok {
					fmt.Println("Calling GetUrl")
					go GetUrl(surl, csite, done)
					numGos += 1
				}
			}
			visited[surl] += 1
		}
	}
}

// Seed starts the crawling process by feeding the URL channel a URL.
func Seed(curl chan []byte) {
	curl <- []byte("https://news.ycombinator.com")
}

// GetUrls parses a site object and looks for links to sites.
func GetUrls(curl chan []byte, site Site, regex *regexp.Regexp) {
	fmt.Println("Making' matches")
	matches := regex.FindAllSubmatch(site.Body, -1)
	for _, match := range matches {
		fmt.Printf("Putting %v on curl\n", match[1])
		curl <- match[1]
	}
}

