package main

import (
	"encoding/json"
	"net/url"
	"time"

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
	contentParsers = append(contentParsers, &JsonParser{})

	crawler := crawler.NewCrawler(initialUrls, contentStorage, &crawler.Config{
		MaxRedirects:    5,
		RevisitDelay:    time.Hour * 2,
		WorkerCount:     100,
		ExcludePatterns: []string{},
	})

	// Adding custom parser to the crawler
	crawler.AddContentParser(&JsonParser{})

	// Adding custom processor to the crawler
	crawler.AddProcessor(&LoggerProcessor{})

	crawler.Start()
}

// Example of custom processor
type LoggerProcessor struct {
}

func (l *LoggerProcessor) Process(result crawler.CrawlResult) error {
	log.Print("Processing result")
	return nil
}

// Example of custom parser
type JsonParser struct {
}

func (p *JsonParser) IsSupportedExtension(extension string) bool {
	for _, supportedMimeTypes := range []string{"application/json"} {
		if extension == supportedMimeTypes {
			return true
		}
	}
	return true
}

func (p *JsonParser) Parse(content string) ([]parser.Token, error) {
	jsonData := map[string]interface{}{}
	err := json.Unmarshal([]byte(content), &jsonData)
	if err != nil {
		return nil, err
	}
	tokens := []parser.Token{}
	for key, value := range jsonData {
		if valueString, ok := value.(string); ok {
			tokens = append(tokens, parser.Token{
				Name:  key,
				Value: valueString,
			})
		}
	}
	return tokens, nil
}
