package crawler

type Processor interface {
	Process(CrawlResult) error
}
