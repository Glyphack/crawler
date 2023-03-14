package crawler

import (
	"net/url"
	"path"

	"github.com/glyphack/crawler/internal/storage"
)

type Processor interface {
	Process(CrawlResult) error
}

func SaveResult(result CrawlResult, storageBackend storage.Storage) error {
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

func getSavePath(url *url.URL) string {
	fileName := url.Path + "-page"
	savePath := path.Join(url.Host, fileName)
	return savePath
}
