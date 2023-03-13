package crawler

import (
	"net/url"
	"path"

	"github.com/glyphack/crawler/internal/parser"
	"github.com/glyphack/crawler/internal/storage"
	log "github.com/sirupsen/logrus"
)

func SaveResult(result WorkerResult, storageBackend storage.Storage) error {
	savePath := getSavePath(result.Url)
	log.Printf("Saving result to %s", savePath)

	switch result.ContentType {
	default:
		err := storageBackend.Set(savePath, string(result.Body))
		if err != nil {
			return err
		}
	}

	return nil
}

func ExtractLinks(content string, p parser.Parser) ([]*url.URL, error) {
	foundUrls := make([]*url.URL, 0)
	parsedUrls, err := p.Parse(content)
	if err != nil {
		log.Printf("Error parsing content: %s", err)
		return foundUrls, err
	}
	log.Printf("Found %d urls", len(parsedUrls))
	for _, parsedUrl := range parsedUrls {
		newUrl, err := url.Parse(parsedUrl.Value)
		if err != nil {
			log.Printf("Error parsing url: %s", err)
			continue
		}
		params := newUrl.Query()
		for param := range params {
			newUrl = stripQueryParam(newUrl, param)
		}
		log.Printf("Url after strip: %s", newUrl)
		if newUrl.Scheme == "http" || newUrl.Scheme == "https" {
			foundUrls = append(foundUrls, newUrl)
		}
	}
	return foundUrls, nil
}

func getSavePath(url *url.URL) string {
	fileName := url.Path
	if fileName == "" {
		fileName = "index.html"
	}
	savePath := path.Join(url.Host, fileName)
	return savePath
}

func stripQueryParam(inputURL *url.URL, stripKey string) *url.URL {
	query := inputURL.Query()
	query.Del(stripKey)
	inputURL.RawQuery = query.Encode()
	return inputURL
}
