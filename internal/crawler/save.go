package crawler

import (
	"net/url"
	"path"

	"github.com/glyphack/crawler/internal/storage"
)

type SaveToFile struct {
	storageBackend storage.Storage
}

func (s *SaveToFile) Process(result CrawlResult) error {
	savePath := getSavePath(result.Url)

	switch result.ContentType {
	default:
		savePath = savePath + ".html"
		err := s.storageBackend.Set(savePath, string(result.Body))
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
