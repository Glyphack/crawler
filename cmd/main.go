package main

import (
	"net/url"

	"github.com/glyphack/crawler/internal/crawler"
	"github.com/glyphack/crawler/internal/parser"
	"github.com/glyphack/crawler/internal/storage"
)

func main() {
	initialUrls := []url.URL{}
	done := make(chan bool)

	myUrl, _ := url.Parse("https://glyphack.com")
	initialUrls = append(initialUrls, *myUrl)

	contentStorage, err := storage.NewFileStorage("./data")
	if err != nil {
		panic(err)
	}

	contentParsers := []parser.Parser{}
	contentParsers = append(contentParsers, &parser.HtmlParser{})

	crawler := crawler.NewCrawler(initialUrls, contentStorage, 10, contentParsers)
	go crawler.Start()
	defer crawler.Terminate()

	<-done
}
