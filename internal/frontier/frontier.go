package frontier

import (
	"net/url"
)

type Frontier struct {
	urls chan *url.URL
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
	f.urls <- url

}

func (f *Frontier) Get() chan *url.URL {
	return f.urls
}
