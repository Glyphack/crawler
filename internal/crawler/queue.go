package crawler

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/glyphack/crawler/internal/frontier"
)

func distributeUrls(frontier *frontier.Frontier, distributedInputs []chan *url.URL) {
	for url := range frontier.Get() {
		index := rand.Intn(len(distributedInputs))
		distributedInputs[index] <- url
	}
}

func mergeResults(workersResults []chan http.Response, out chan http.Response) {
	collect := func(in chan http.Response) {
		for result := range in {
			log.Printf("Got result to collect %s", result.Request.URL)
			// Handle closed
			out <- result
			log.Printf("Collected result %s", result.Request.URL)
		}
		log.Println("Worker finished")
	}

	for i, result := range workersResults {
		log.Printf("Start collecting results from worker %d", i)
		go collect(result)
	}
}
