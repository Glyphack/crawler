package crawler

import (
	"log"
	"net/http"
	"net/url"

	"github.com/glyphack/crawler/internal/fetch"
)

type Worker struct {
	input      chan *url.URL
	deadLetter chan *url.URL
	result     chan http.Response
	done       chan struct{}
	id         int
}

func NewWorker(input chan *url.URL, result chan http.Response, done chan struct{}, id int) *Worker {
	return &Worker{
		input:  input,
		result: result,
		done:   done,
		id:     id,
	}
}

func (w *Worker) Start() {
	log.Printf("Worker %d started", w.id)
	for {
		select {
		// TODO Handle failures and retry
		// case failedUrl := <-w.deadLetter:
		case url := <-w.input:
			log.Printf("Worker %d fetching %s", w.id, url)
			content, err := fetch.Fetch(url)
			if err != nil {
				w.deadLetter <- url
				log.Printf("Worker %d error fetching content: %s", err)
				continue
			}
			log.Printf("Worker %d fetched %s", w.id, url)
			w.result <- *content
		case <-w.done:
			return
		}
	}
}
