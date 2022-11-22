package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var (
	config = &tls.Config{
		InsecureSkipVerify: false,
	}
	transport = &http.Transport{
		TLSClientConfig: config,
	}
	newClient = &http.Client{
		Transport: transport,
	}
	queue      = make(chan string)
	hasVisited = make(map[string]bool)
)

func main() {
	arguments := os.Args[1:]
	if len(arguments) == 0 {
		fmt.Println("Missing URL")
		os.Exit(1)
	}
	//baseURL := "https://www.youtube.com/feed/library"
	baseURL := arguments[0]
	go func() {
		queue <- baseURL
	}()
	for href := range queue {
		if !hasVisited[href] && isSameDomain(href, baseURL) {
			crawlURL(href)
		}
	}
}

func crawlURL(href string) {
	hasVisited[href] = true
	fmt.Printf("Crawling URL-->%v\n", href)
	response, err := newClient.Get(href)
	defer response.Body.Close()
	checkErr(err)
	links, err := extractlinks.All(response.Body)
	checkErr(err)
	for _, link := range links {
		absoluteURL := toFixedURL(link.Href, href)
		go func() {
			queue <- absoluteURL
		}()
	}
}

func toFixedURL(href string, baseURL string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	toFixedURI := base.ResolveReference(uri)
	return toFixedURI.String()
}

func isSameDomain(href string, baseURL string) bool {
	uri, err := url.Parse(href)
	if err != nil {
		return false
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return false
	}
	if uri.Host != base.Host {
		return false
	}
	return true
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
