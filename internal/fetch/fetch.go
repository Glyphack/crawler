package fetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func Fetch(url *url.URL) (string, error) {
	res, err := http.Get(url.String())
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
