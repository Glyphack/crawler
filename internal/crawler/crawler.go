package crawler

import (
	"net/url"

	log "github.com/sirupsen/logrus"

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
	deadLetter      chan *url.URL
}

func NewCrawler(initialUrls []url.URL, contentStorage storage.Storage, workerCount int, contentParser []parser.Parser) *Crawler {
	deadLetter := make(chan *url.URL)
	return &Crawler{
		frontier:      frontier.NewFrontier(initialUrls),
		storage:       contentStorage,
		workerCount:   workerCount,
		contentParser: contentParser,
		deadLetter:    deadLetter,
	}
}

func (c *Crawler) Start() {
	distributedInputs := make([]chan *url.URL, c.workerCount)
	workersResults := make([]chan WorkerResult, c.workerCount)
	done := make(chan struct{})

	for i := 0; i < c.workerCount; i++ {
		distributedInputs[i] = make(chan *url.URL)
		workersResults[i] = make(chan WorkerResult)
	}
	go distributeUrls(c.frontier, distributedInputs)
	for i := 0; i < c.workerCount; i++ {
		worker := NewWorker(distributedInputs[i], workersResults[i], done, i, c.deadLetter)
		go worker.Start()
	}

	mergedResults := make(chan WorkerResult)
	go mergeResults(workersResults, mergedResults)
	// processedSignal := make(chan struct{}, c.workerCount)
	newUrls := make(chan *url.URL)
	go func() {
		for newUrl := range newUrls {
			_ = c.frontier.Add(newUrl)
			// if !added {
			// 	log.Warnf("Url %s already in frontier", newUrl)
			// 	go func() {
			// 		processedSignal <- struct{}{}
			// 	}()
			// }
		}
	}()

	go func() {
		for deadUrl := range c.deadLetter {
			log.Debugf("Dismissed %s", deadUrl)
		}
	}()

	for result := range mergedResults {
		err := SaveResult(result, c.storage)
		if err != nil {
			log.Error(err)
		}
		links, err := ExtractLinks(string(result.Body), c.contentParser[0])
		if err != nil {
			log.Error(err)
		}
		go func() {
			for _, link := range links {
				newUrls <- link
			}
		}()
	}
	log.Println("Crawler exited")
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
