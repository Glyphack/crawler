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
	exclude     []string
}

func NewFrontier(initialUrls []url.URL, exclude []string) *Frontier {
	history := make(map[url.URL]time.Time)
	f := &Frontier{
		urls:    make(chan *url.URL, len(initialUrls)),
		history: history,
		exclude: exclude,
	}

	for _, u := range initialUrls {
		f.Add(&u)
	}
	return f
}

func (f *Frontier) Add(url *url.URL) bool {
	if f.terminating {
		return false
	}
	if f.Seen(url) {
		log.WithFields(log.Fields{
			"url": url,
		}).Info("Already seen")
		return false
	}
	for _, pattern := range f.exclude {
		if pattern == url.Host {
			log.WithFields(log.Fields{
				"url": url,
			}).Info("Excluded")
			return false
		}
	}
	f.history[*url] = time.Now()
	f.urls <- url

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
