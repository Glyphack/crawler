package parser

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type HtmlParser struct {
}

func (p *HtmlParser) getSupportedExtensions() []string {
	return []string{".html", ".htm"}
}

func (p *HtmlParser) IsSupportedExtension(extension string) bool {
	for _, supportedExtension := range p.getSupportedExtensions() {
		if extension == supportedExtension {
			return true
		}
	}
	return true
}

func (p *HtmlParser) Parse(content string) ([]Token, error) {
	htmlParser := html.NewTokenizer(strings.NewReader(content))
	tokens := []Token{}
	for {
		tokenType := htmlParser.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := htmlParser.Token()
		if tokenType == html.StartTagToken {
			switch token.Data {
			case "a":
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						tokens = append(tokens, Token{Name: "link", Value: attr.Val})
					}
				}
			}
		}
	}

	if htmlParser.Err() != nil {
		if !errors.Is(htmlParser.Err(), io.EOF) {
			return tokens, fmt.Errorf("error scanning html: %s", htmlParser.Err())
		}
	}
	return tokens, nil
}
