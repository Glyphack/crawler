package fetch

import (
	"fmt"
	"net/http"
	"net/url"
)

func Fetch(url *url.URL) (*http.Response, error) {
	res, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return res, nil
}
