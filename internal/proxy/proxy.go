package proxy

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	proxies []string
	mu      sync.Mutex
	index   int
}

var (
	clientTimeout = 5 * time.Second
)

func normalizeProxyURL(raw string) string {
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		return "http://" + raw
	}
	return raw
}


func NewManager(proxyFile string) (*Manager, error) {
	proxies, err := loadProxies(proxyFile)
	if err != nil {
		return nil, err
	}

	validated := validateProxies(proxies)
	if len(validated) == 0 {
		return nil, errors.New("no valid proxies found")
	}

	log.Printf("[proxy] %d proxies loaded and validated\n", len(validated))
	return &Manager{proxies: validated, index: 0}, nil
}

func loadProxies(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var proxies []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			proxies = append(proxies, line)
		}
	}
	return proxies, nil
}

func validateProxies(proxies []string) []string {
	var valid []string
	for _, p := range proxies {
		if testProxy(p) {
			valid = append(valid, p)
		}
	}
	return valid
}

func testProxy(proxyURL string) bool {
	normalized := normalizeProxyURL(proxyURL)
	fmt.Println("Testing proxy:", normalized)
	proxy, _ := url.Parse(normalized)
	transport := &http.Transport{Proxy: http.ProxyURL(proxy)}
	client := &http.Client{
		Transport: transport,
		Timeout:   clientTimeout,
	}

	req, _ := http.NewRequest("GET", "https://www.instagram.com", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	resp.Body.Close()
	return true
}

// GetNext returns the next proxy in round-robin
func (pm *Manager) GetNext() string {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	proxy := pm.proxies[pm.index]
	pm.index = (pm.index + 1) % len(pm.proxies)
	return proxy
}

// GetRandom returns a random proxy
func (pm *Manager) GetRandom() string {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.proxies[rand.Intn(len(pm.proxies))]
}
