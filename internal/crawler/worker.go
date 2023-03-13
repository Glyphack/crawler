package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

type WorkerResult struct {
	Url         *url.URL
	ContentType string
	Body        []byte
}

type Worker struct {
	input      chan *url.URL
	deadLetter chan *url.URL
	result     chan WorkerResult
	done       chan struct{}
	id         int

	// Only contains the host part of the URL
	history map[string]time.Time
}

func NewWorker(input chan *url.URL, result chan WorkerResult, done chan struct{}, id int, deadLetter chan *url.URL) *Worker {
	history := make(map[string]time.Time)
	return &Worker{
		input:      input,
		result:     result,
		done:       done,
		id:         id,
		history:    history,
		deadLetter: deadLetter,
	}
}

func (w *Worker) Start() {
	log.Printf("Worker %d started", w.id)
	for {
		select {
		case url := <-w.input:
			content, err := w.fetch(url)
			if err != nil {
				log.Printf("Worker %d error fetching content: %s", w.id, err)
				w.deadLetter <- url
				continue
			}
			w.result <- content
		case deadUrl := <-w.deadLetter:
			content, err := w.fetch(deadUrl)
			if err != nil {
				log.Printf("Worker %d error fetching content: %s", w.id, err)
				w.deadLetter <- deadUrl
				continue
			}
			w.result <- content
		case <-w.done:
			return
		}
	}
}

func (w *Worker) CheckPoliteness(url *url.URL) bool {
	if lastFetch, ok := w.history[url.Host]; ok {
		return time.Since(lastFetch) > 2*time.Second
	}
	return true
}

func (w *Worker) fetch(url *url.URL) (WorkerResult, error) {
	log.Printf("Worker %d fetching %s", w.id, url)
	defer log.Printf("Worker %d done fetching %s", w.id, url)
	defer func() {
		w.history[url.Host] = time.Now()
	}()
	for !w.CheckPoliteness(url) {
		log.Printf("Worker %d waiting for %s", w.id, url)
		time.Sleep(2 * time.Second)
	}
	res, err := http.Get(url.String())
	if err != nil {
		return WorkerResult{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return WorkerResult{}, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return WorkerResult{}, err
	}

	var inferredContentType string
	contentType, ok := res.Header["Content-Type"]
	if ok && len(contentType) > 0 {
		inferredContentType = contentType[0]
	} else {
		inferredContentType = http.DetectContentType(body)
	}

	return WorkerResult{
		Url:         url,
		ContentType: inferredContentType,
		Body:        body,
	}, nil
}
