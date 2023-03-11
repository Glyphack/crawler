package crawler

import (
	"log"
	"net/url"
	"path"

	"github.com/glyphack/crawler/internal/fetch"
	"github.com/glyphack/crawler/internal/frontier"
	"github.com/glyphack/crawler/internal/storage"
)

type Crawler struct {
	excludePatterns []string
	frontier        *frontier.Frontier
	storage         storage.Storage
}

func NewCrawler(initialUrls []url.URL) *Crawler {
	return &Crawler{
		frontier: frontier.NewFrontier(initialUrls),
		storage:  storage.NewFileStorage("."),
	}
}

func (c *Crawler) Start() {
	url := <-c.frontier.Get()
	page, err := fetch.Fetch(url)
	if err != nil {
		log.Printf("Error fetching %s: %s", url, err)
		return
	}

	savePath := path.Join(".", url.Host, url.Path)
	err = c.storage.Set(savePath, page)

	if err != nil {
		log.Printf("Error saving %s: %s", savePath, err)
		return
	}
}
