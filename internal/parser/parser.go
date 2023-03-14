package parser

type Token struct {
	Name  string
	Value string
}

type Parser interface {
	IsSupportedExtension(extension string) bool
	Parse(content string) ([]Token, error)
}
