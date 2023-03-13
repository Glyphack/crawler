package crawler

import (
	"fmt"
	"net/url"
	"path"

	"github.com/glyphack/crawler/internal/parser"
	"github.com/glyphack/crawler/internal/storage"
	log "github.com/sirupsen/logrus"
)

func SaveResult(result WorkerResult, storageBackend storage.Storage) error {
	savePath := getSavePath(result.Url)

	switch result.ContentType {
	default:
		savePath = savePath + ".html"
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
		return foundUrls, fmt.Errorf("Error parsing content: %s", err)
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
	return foundUrls, nil
}

func getSavePath(url *url.URL) string {
	fileName := url.Path + "-page"
	savePath := path.Join(url.Host, fileName)
	return savePath
}

func stripQueryParam(inputURL *url.URL, stripKey string) *url.URL {
	query := inputURL.Query()
	query.Del(stripKey)
	inputURL.RawQuery = query.Encode()
	return inputURL
}
