package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Site struct {
	Url, Body []byte
}

func getUrl(url []byte, csite chan Site, death chan string) {
	resource := string(url)
	defer func () {
		death <- resource
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


func main() {
	regex := regexp.MustCompile("<a.*?href=[\"'](http.*?)[\"']")

	// URL seeds
	curl := make(chan []byte)
	go func() {
		curl <- []byte("http://joearms.github.com/2013/03/28/solving-the-wrong-problem.html")
		curl <- []byte("http://www.forbes.com/fdc/welcome_mjx.shtml")
		curl <- []byte("https://segment.io/academy/email-is-the-easiest-way-to-improve-retention")
	}()

	csite := make(chan Site)
	death := make(chan string)

	go func () {
		// Ensure we don't have many open connections
		visited := make(map[string]int)
		i := 0
		for {
			if i > 10 {
				<-death
				i -= 1
			}
			url := string(<-curl)
			if _, ok := visited[url]; !ok {
				go getUrl([]byte(url), csite, death)
				i += 1
			}
			visited[url] += 1
		}
	}()

	for {
		site := <-csite
		go func () {
			matches := regex.FindAllSubmatch(site.Body, -1)
			for _, match := range matches {
				curl <- match[1]
			}
		}()
	}
}
