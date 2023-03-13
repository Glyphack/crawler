package frontier

import (
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

type Frontier struct {
	urls        chan *url.URL
	terminating bool
	history     map[url.URL]time.Time
}

func NewFrontier(initialUrls []url.URL) *Frontier {
	history := make(map[url.URL]time.Time)
	f := &Frontier{
		urls:    make(chan *url.URL, len(initialUrls)),
		history: history,
	}

	for _, u := range initialUrls {
		f.Add(&u)
	}
	return f
}

func (f *Frontier) Add(url *url.URL) bool {
	log.Printf("Adding %s to frontier", url)
	if f.terminating {
		return false
	}
	if f.Seen(url) {
		return false
	}
	f.history[*url] = time.Now()
	f.urls <- url

	log.Printf("Added %s to frontier", url)
	return true
}

func (f *Frontier) Get() chan *url.URL {
	return f.urls
}

func (f *Frontier) Terminate() {
	close(f.urls)
	f.terminating = true
}

func (f *Frontier) Seen(url *url.URL) bool {
	if lastFetch, ok := f.history[*url]; ok {
		return time.Since(lastFetch) < 2*time.Hour
	}
	return false
}
