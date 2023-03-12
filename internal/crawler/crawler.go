package crawler

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/glyphack/crawler/internal/frontier"
	"github.com/glyphack/crawler/internal/parser"
	"github.com/glyphack/crawler/internal/storage"
)

type Crawler struct {
	excludePatterns []string
	frontier        *frontier.Frontier
	storage         storage.Storage
	contentParser   []parser.Parser
	workerCount     int
	deadLetter      chan http.Response
}

func NewCrawler(initialUrls []url.URL, contentStorage storage.Storage, workerCount int, contentParser []parser.Parser) *Crawler {
	return &Crawler{
		frontier:      frontier.NewFrontier(initialUrls),
		storage:       contentStorage,
		workerCount:   workerCount,
		contentParser: contentParser,
	}
}

func (c *Crawler) Start() {
	distributedInputs := make([]chan *url.URL, c.workerCount)
	workersResults := make([]chan http.Response, c.workerCount)
	done := make(chan struct{})

	for i := 0; i < c.workerCount; i++ {
		distributedInputs[i] = make(chan *url.URL)
		workersResults[i] = make(chan http.Response)
	}

	go distributeUrls(c.frontier, distributedInputs)

	for i := 0; i < c.workerCount; i++ {
		worker := NewWorker(distributedInputs[i], workersResults[i], done, i)
		go worker.Start()
	}

	mergedResults := make(chan http.Response, 100)
	go mergeResults(workersResults, mergedResults)
	for result := range mergedResults {
		log.Printf("Got result for %s", result.Request.URL)
		resultBody, err := io.ReadAll(result.Body)
		result.Body.Close()
		if err != nil {
			log.Printf("Error reading content: %s", err)
			// c.deadLetter <- result
			continue
		}
		savePath := path.Join(result.Request.Host, result.Request.URL.Path)
		err = c.storage.Set(savePath, string(resultBody))
		log.Printf("Saved content to %s", savePath)
		if err != nil {
			log.Printf("Error saving content: %s", err)
			// c.deadLetter <- result
			continue
		}

		parserHandled := false
		for _, parser := range c.contentParser {
			if parser.IsSupportedExtension(result.Request.URL.Path) {
				parserHandled = true
				log.Printf("Parsing content from %s with %s", result.Request.URL, c.contentParser)
				parsedUrls, err := parser.Parse(string(resultBody))
				if err != nil {
					log.Printf("Error parsing content: %s", err)
					continue
				}
				for _, parsedUrl := range parsedUrls {
					url, err := url.Parse(parsedUrl.Value)
					if err != nil {
						log.Printf("Error parsing url: %s", err)
						continue
					}
					if url.Scheme == "http" || url.Scheme == "https" {
						c.frontier.Add(url)
					}
				}
			}
			if !parserHandled {
				log.Printf("No parser found for %s", result.Request.URL)
			}
		}
		log.Printf("Finished processing %s", result.Request.URL)
	}

	log.Println("Crawler finished")
}

func (c *Crawler) Terminate() {
	c.frontier.Terminate()
}

func (c *Crawler) AddContentParser(contentParser parser.Parser) {
	c.contentParser = append(c.contentParser, contentParser)
}

func (c *Crawler) AddExcludePattern(pattern string) {
	c.excludePatterns = append(c.excludePatterns, pattern)
}
