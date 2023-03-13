package fetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Fetch(url *url.URL) ([]byte, error) {
	res, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	contentType, ok := res.Header["Content-Type"]

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
