package crawler

import (
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/glyphack/crawler/internal/frontier"
	"github.com/glyphack/crawler/internal/parser"
	"github.com/glyphack/crawler/internal/storage"
)

type Crawler struct {
	config         *Config
	frontier       *frontier.Frontier
	storage        storage.Storage
	contentParsers []parser.Parser
	deadLetter     chan *url.URL
	processors     []Processor
}

func NewCrawler(initialUrls []url.URL, contentStorage storage.Storage, config *Config) *Crawler {
	deadLetter := make(chan *url.URL)
	contentParser := []parser.Parser{&parser.HtmlParser{}}
	return &Crawler{
		frontier:       frontier.NewFrontier(initialUrls, config.ExcludePatterns),
		storage:        contentStorage,
		contentParsers: contentParser,
		deadLetter:     deadLetter,
		config:         config,
	}
}

func (c *Crawler) Start() {
	distributedInputs := make([]chan *url.URL, c.config.WorkerCount)
	workersResults := make([]chan CrawlResult, c.config.WorkerCount)
	done := make(chan struct{})

	for i := 0; i < c.config.WorkerCount; i++ {
		distributedInputs[i] = make(chan *url.URL)
		workersResults[i] = make(chan CrawlResult)
	}
	go distributeUrls(c.frontier, distributedInputs)
	for i := 0; i < c.config.WorkerCount; i++ {
		worker := NewWorker(distributedInputs[i], workersResults[i], done, i, c.deadLetter)
		go worker.Start()
	}

	mergedResults := make(chan CrawlResult)
	go mergeResults(workersResults, mergedResults)
	newUrls := make(chan *url.URL)
	c.AddProcessor(&LinkExtractor{Parsers: c.contentParsers, NewUrls: newUrls})
	c.AddProcessor(&SaveToFile{storageBackend: c.storage})
	go func() {
		for newUrl := range newUrls {
			_ = c.frontier.Add(newUrl)
		}
	}()

	go func() {
		for deadUrl := range c.deadLetter {
			log.Debugf("Dismissed %s", deadUrl)
		}
	}()

	for result := range mergedResults {
		for _, processor := range c.processors {
			go func(processor Processor, result CrawlResult) {
				processErr := processor.Process(result)
				if processErr != nil {
					log.Error(processErr)
				}
			}(processor, result)
		}
	}
	log.Println("Crawler exited")
}

func (c *Crawler) Terminate() {
	c.frontier.Terminate()
}
func (c *Crawler) AddContentParser(contentParser parser.Parser) {
	c.contentParsers = append(c.contentParsers, contentParser)
}

func (c *Crawler) AddExcludePattern(pattern string) {
	c.config.ExcludePatterns = append(c.config.ExcludePatterns, pattern)
}

func (c *Crawler) AddProcessor(processor Processor) {
	c.processors = append(c.processors, processor)
}
