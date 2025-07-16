package fetcher

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/0xTan1319/instagram-scraper/internal/config"
	"github.com/0xTan1319/instagram-scraper/internal/proxy"
)

type Fetcher struct {
	ProxyManager *proxy.Manager
	UserAgent    string
	Timeout      time.Duration
	MaxRetries   int
}

// NewFetcher creates a fetcher with retry & proxy support
func NewFetcher(pm *proxy.Manager) *Fetcher {
	return &Fetcher{
		ProxyManager: pm,
		UserAgent:    config.AppConfig.UserAgent,
		Timeout:      config.AppConfig.RequestTimeout,
		MaxRetries:   5,
	}
}

// Get performs a GET request to the given URL using rotating proxies
func (f *Fetcher) Get(targetURL string) ([]byte, error) {
	var lastErr error

	for i := 0; i < f.MaxRetries; i++ {
		proxyStr := f.ProxyManager.GetNext()

		client, err := f.buildClientWithProxy(proxyStr)
		if err != nil {
			log.Printf("[fetcher] failed to build proxy client: %v", err)
			continue
		}

		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", f.UserAgent)
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[fetcher] request failed (attempt %d): %v", i+1, err)
			lastErr = err
			continue
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("[fetcher] non-200 response (attempt %d): %d", i+1, resp.StatusCode)
			lastErr = errors.New("non-200 response")
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	}

	return nil, lastErr
}

func (f *Fetcher) buildClientWithProxy(proxyStr string) (*http.Client, error) {
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	return &http.Client{
		Timeout:   f.Timeout,
		Transport: transport,
	}, nil
}
