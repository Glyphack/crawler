package frontier

import (
	"log"
	"net/url"
)

type Frontier struct {
	urls        chan *url.URL
	terminating bool
}

func NewFrontier(initialUrls []url.URL) *Frontier {
	f := &Frontier{
		urls: make(chan *url.URL, len(initialUrls)),
	}

	for _, u := range initialUrls {
		f.Add(&u)
	}
	return f
}

func (f *Frontier) Add(url *url.URL) {
	log.Printf("Adding %s to frontier", url)
	if f.terminating {
		return
	}
	f.urls <- url

}

func (f *Frontier) Get() chan *url.URL {
	return f.urls
}

func (f *Frontier) Terminate() {
	close(f.urls)
	f.terminating = true
}
