package main

import (
	"fmt"
	"log"

	"github.com/0xTan1319/instagram-scraper/internal/config"
	"github.com/0xTan1319/instagram-scraper/internal/fetcher"
	"github.com/0xTan1319/instagram-scraper/internal/parser"
	"github.com/0xTan1319/instagram-scraper/internal/proxy"
)

func main() {
	config.LoadConfig()

	log.Println("🚀 Starting Instagram Scraper...")

	pm, err := proxy.NewManager(config.AppConfig.ProxyFile)
	if err != nil {
		log.Fatalf("❌ Failed to initialize proxy manager: %v", err)
	}

	currentProxy := pm.GetNext()
	if currentProxy == "" {
		log.Fatalf("❌ No valid proxy available")
	}

	log.Println("🌐 Using proxy:", currentProxy)

	f := fetcher.NewFetcher(pm)

	html, err := f.Get("https://www.instagram.com/explore/tags/fashion/")
	if err != nil {
		log.Fatalf("❌ Failed to fetch page: %v", err)
	}

	// Safely limit output to 500 characters
	previewLength := min(len(html), 500)
	fmt.Println(string(html[:previewLength]))
	fmt.Println("📄 Page Preview:\n", string(html[:previewLength]))

	videos, err := parser.ParseVideos(string(html))
	if err != nil {
		log.Fatalf("❌ Failed to parse videos: %v", err)
	}

	for _, v := range videos {
		fmt.Println("🎬 Video:", v.VideoURL)
		fmt.Println("📸 Thumbnail:", v.Thumbnail)
		fmt.Println("📝 Caption:", v.Caption)
		fmt.Println("👤 Author:", v.Author)
		fmt.Println("---")
	}
}
