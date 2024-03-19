package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	extractlinks "github.com/hvullas/goWebCrawler/extractLinks"
)

var (
	config = &tls.Config{
		InsecureSkipVerify: true,
	}
	transport = &http.Transport{
		TLSClientConfig: config,
	}

	netclient = &http.Client{
		Transport: transport,
	}

	queue      = make(chan string)
	hasVisited = make(map[string]bool)
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Println("Missing URL")
		os.Exit(1)
	}

	go func() {
		queue <- args[0]
	}()

	for href := range queue {
		if !hasVisited[href] {
			crawlURL(href)
		}

	}

}

func crawlURL(href string) {
	hasVisited[href] = true
	fmt.Printf("Crawling url -> %v \n", href)

	response, err := netclient.Get(href)
	checkErr(err)
	defer response.Body.Close()

	links, err := extractlinks.All(response.Body)
	checkErr(err)

	for _, link := range links {
		absURL := tofixedURL(link.Href, href)
		go func() {
			queue <- absURL
		}()

	}

}

func tofixedURL(href, baseURL string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	tofixedUri := base.ResolveReference(uri)
	return tofixedUri.String()
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
