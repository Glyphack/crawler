package main

import (
	"net/url"

	"github.com/glyphack/crawler/internal/crawler"
)

func main() {
	initialUrls := []url.URL{}

	myUrl, _ := url.Parse("https://glyphack.com")
	initialUrls = append(initialUrls, *myUrl)

	// Create a new crawler
	crawler := crawler.NewCrawler(initialUrls)

	crawler.Start()
}
