# Crawler

A web crawler written in golang.

## Usage

Set inial Urls & storage config then run the crawler.

```go
func main() {
    initialUrls := []url.URL{}

    myUrl, _ := url.Parse("https://glyphack.com")
    initialUrls = append(initialUrls, *myUrl)

    contentStorage, err := storage.NewFileStorage("./data")
    if err != nil {
        panic(err)
    }

    contentParsers := []parser.Parser{}
    contentParsers = append(contentParsers, &parser.HtmlParser{})

    crawler := crawler.NewCrawler(initialUrls, contentStorage, 10, contentParsers)
    crawler.Start()
}
```

## Config

The following options are supported:

- maxRedirects: Number of maximum redirects to follow
- RevisitCoolDown: Number of seconds to wait before revisiting a URL
- workerCount: Number of workers that simultaneously visit URLs

## Extensibility

You can extend the the crawler by adding new storage to it.

Each visited website will be handled by a content parser.
You can implement custom parsers for different MIME types e.g. `application/pdf`.
