package crawler

import (
	"math/rand"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/glyphack/crawler/internal/frontier"
)

func distributeUrls(frontier *frontier.Frontier, distributedInputs []chan *url.URL) {
	for url := range frontier.Get() {
		index := rand.Intn(len(distributedInputs))
		distributedInputs[index] <- url
		log.Printf("Distributed %s to worker %d", url, index)
	}
}

func mergeResults(workersResults []chan WorkerResult, out chan WorkerResult) {
	collect := func(in chan WorkerResult) {
		for result := range in {
			log.Printf("Got result to collect %s", result.Url)
			out <- result
			log.Printf("Collected result %s", result.Url)
		}
		log.Println("Worker finished")
	}

	for i, result := range workersResults {
		log.Printf("Start collecting results from worker %d", i)
		go collect(result)
	}
}
