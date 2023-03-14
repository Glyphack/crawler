package crawler

import "github.com/glyphack/crawler/internal/storage"

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
