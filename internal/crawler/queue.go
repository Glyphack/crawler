package crawler

import (
	"math/rand"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/glyphack/crawler/internal/frontier"
)

func distributeUrls(frontier *frontier.Frontier, distributedInputs []chan *url.URL) {
	HostToWorker := make(map[string]int)
	for url := range frontier.Get() {
		index := rand.Intn(len(distributedInputs))
		if prevIndex, ok := HostToWorker[url.Host]; ok {
			index = prevIndex
		} else {
			HostToWorker[url.Host] = index
		}
		distributedInputs[index] <- url
	}
}

func mergeResults(workersResults []chan CrawlResult, out chan CrawlResult) {
	collect := func(in chan CrawlResult) {
		for result := range in {
			out <- result
		}
		log.Println("Worker finished")
	}

	for i, result := range workersResults {
		log.Printf("Start collecting results from worker %d", i)
		go collect(result)
	}
}
