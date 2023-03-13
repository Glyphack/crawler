package crawler

import (
	"fmt"
	"net/url"

	"github.com/glyphack/crawler/internal/parser"
	log "github.com/sirupsen/logrus"
)

type LinkExtractor struct {
	Parser  parser.Parser
	NewUrls chan *url.URL
}

func (e *LinkExtractor) Process(result WorkerResult) error {
	foundUrls := make([]*url.URL, 0)
	parsedUrls, err := e.Parser.Parse(string(result.Body))
	if err != nil {
		return fmt.Errorf("Error parsing content: %s", err)
	}
	log.Infof("Extracted %d urls", len(parsedUrls))
	for _, parsedUrl := range parsedUrls {
		newUrl, err := url.Parse(parsedUrl.Value)
		if err != nil {
			log.Debugf("Error parsing url: %s", err)
			continue
		}
		params := newUrl.Query()
		for param := range params {
			newUrl = stripQueryParam(newUrl, param)
		}
		if newUrl.Scheme == "http" || newUrl.Scheme == "https" {
			foundUrls = append(foundUrls, newUrl)
		}
	}
	for _, foundUrl := range foundUrls {
		e.NewUrls <- foundUrl
	}
	return nil
}

func stripQueryParam(inputURL *url.URL, stripKey string) *url.URL {
	query := inputURL.Query()
	query.Del(stripKey)
	inputURL.RawQuery = query.Encode()
	return inputURL
}
