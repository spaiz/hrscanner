package main

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"fmt"
)

var httpClient *http.Client

// Client responsible for doing http request to specific host,
// and returning headers
type Client struct {
	DNSResolver DNSResolver
}

func init() {
	httpClient = createHTTPClient()
}

// it's safe to use single http client in multiple goroutines,
// so we create it once and reuse
func createHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost: 1,
		DisableKeepAlives:   true,
	}

	return &http.Client{
		Transport: tr,
		Timeout:   time.Second * 60,
	}
}

func (r *Client) getHeaders(host string) (http.Header, error) {
	ips, err := r.DNSResolver.Resolve(host)
	if err != nil {
		return nil, err
	}

	// Ideally we should loop try all received IPs until we succeed
	// but it's good enough for me to try the first one only
	req, err := http.NewRequest("HEAD", fmt.Sprintf("http://%s", ips[0]), nil)
	req.Host = host
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	io.Copy(ioutil.Discard, res.Body)
	defer res.Body.Close()

	return res.Header, err
}

// GetHeaders returns response headers.
// it's separated to be able to wrap it later with retries
func (r *Client) GetHeaders(host string) (http.Header, error) {
	return r.getHeaders(host)
}
