# Crawler

A web crawler written in golang.

## Usage

Set initial Urls & storage config then run the crawler.

```go
func main() {
    log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
    initialUrls := []url.URL{}

    myUrl, _ := url.Parse("https://glyphack.com")
    initialUrls = append(initialUrls, *myUrl)

    contentStorage, err := storage.NewFileStorage("./data")
    if err != nil {
        panic(err)
    }

    contentParsers := []parser.Parser{}
    contentParsers = append(contentParsers, &JsonParser{})

    crawler := crawler.NewCrawler(initialUrls, contentStorage, &crawler.Config{
        MaxRedirects:    5,
        RevisitDelay:    time.Hour * 2,
        WorkerCount:     100,
        ExcludePatterns: []string{},
    })

    // Adding custom parser to the crawler
    crawler.AddContentParser(&JsonParser{})

    // Adding custom processor to the crawler
    crawler.AddProcessor(&LoggerProcessor{})

    crawler.Start()
}
}
```

## Config

The following options are supported:

- maxRedirects: Number of maximum redirects to follow
- RevisitCoolDown: Number of seconds to wait before revisiting a URL
- workerCount: Number of workers that simultaneously visit URLs

## Extensibility

You can extend the the crawler by adding new storage to it.

Each visited website will be handled by a content parser. if a parser
You can implement custom parsers for different MIME types e.g. `application/pdf`.

### Parsers

Parsers are used to parse content of web pages to extract links.
The following parsers exist internally:

- html

You can add custom parsers which implement the Parser interface:

```go
type JsonParser struct {
}

func (p *JsonParser) IsSupportedExtension(extension string) bool {
    for _, supportedMimeTypes := range []string{"application/json"} {
        if extension == supportedMimeTypes {
            return true
        }
    }
    return true
}

func (p *JsonParser) Parse(content string) ([]parser.Token, error) {
    jsonData := map[string]interface{}{}
    err := json.Unmarshal([]byte(content), &jsonData)
    if err != nil {
        return nil, err
    }
    tokens := []parser.Token{}
    for key, value := range jsonData {
        if valueString, ok := value.(string); ok {
            tokens = append(tokens, parser.Token{
                Name:  key,
                Value: valueString,
            })
        }
    }
    return tokens, nil
}
```

### Processor

after crawl finishes the application calls registered processors.

Processors can do anything with the result of a webpage,
for example saving content is handled by an internal processor.

**Note**: processors are ran in a separate goroutine.
So sharing a memory can cause data races.

```go
type LoggerProcessor struct {
}

func (l *LoggerProcessor) Process(result crawler.CrawlResult) error {
    log.Print("Processing result")
    return nil
}
```
