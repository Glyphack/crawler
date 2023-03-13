package main

import (
	"net/url"

	"github.com/glyphack/crawler/internal/crawler"
	"github.com/glyphack/crawler/internal/parser"
	"github.com/glyphack/crawler/internal/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	initialUrls := []url.URL{}

	myUrl, _ := url.Parse("https://glyphack.com")
	initialUrls = append(initialUrls, *myUrl)

	contentStorage, err := storage.NewFileStorage("./data")
	if err != nil {
		panic(err)
	}

	contentParsers := []parser.Parser{}
	contentParsers = append(contentParsers, &parser.HtmlParser{})

	crawler := crawler.NewCrawler(initialUrls, contentStorage, 1000, contentParsers)
	crawler.Start()
}
