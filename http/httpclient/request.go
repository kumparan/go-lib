package httpclient

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func ParseURL(rawUrl string, keyVal ...string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	keyValLength := len(keyVal)
	if keyValLength%2 != 0 {
		return "", errors.New("Key-value param must be in key-value format")
	}

	q := u.Query()
	for counter := 0; counter < keyValLength; counter += 2 {
		q.Add(keyVal[counter], keyVal[counter+1])
	}
	finalURL := fmt.Sprintf("%s?%s", u.String(), q.Encode())
	return finalURL, nil
}

func NewRequestWithHostHeader(method, urlStr, hostHeader string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if hostHeader != "" {
		req.Header.Add("Host", hostHeader)
	}
	return req, err
}
