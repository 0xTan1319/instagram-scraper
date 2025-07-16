package parser

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type VideoMeta struct {
	PostURL     string
	VideoURL    string
	Thumbnail   string
	Caption     string
	Author      string
	IsReel      bool
	Timestamp   int64
}

// ParseVideos parses HTML and extracts video posts
func ParseVideos(html string) ([]VideoMeta, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	// Instagram embeds data in <script type="application/ld+json"> or window._sharedData
	scripts := doc.Find("script")

	var jsonRaw string

	scripts.EachWithBreak(func(i int, s *goquery.Selection) bool {
		html, _ := s.Html()
		if strings.Contains(html, "window._sharedData") {
			jsonRaw = html
			return false
		}
		return true
	})

	if jsonRaw == "" {
		return nil, errors.New("no JSON blob found in HTML")
	}

	// Clean the blob
	start := strings.Index(jsonRaw, "{")
	end := strings.LastIndex(jsonRaw, "};")
	if start == -1 || end == -1 {
		return nil, errors.New("invalid JSON blob")
	}

	jsonBlob := jsonRaw[start : end+1]

	// ✅ To parse this JSON, we can use struct or dynamic parsing.
	// For now, simple regex extract video URLs (for proof of concept)
	videoMetas := extractVideoURLs(jsonBlob)
	return videoMetas, nil
}

// Naive regex parser to extract .mp4 video URLs and thumbnails
func extractVideoURLs(json string) []VideoMeta {
	videoURLPattern := regexp.MustCompile(`"video_url":"(https:\\/\\/[^"]+\.mp4)"`)
	thumbPattern := regexp.MustCompile(`"display_url":"(https:\\/\\/[^"]+)"`)
	captionPattern := regexp.MustCompile(`"edge_media_to_caption":{"edges":\[\{"node":\{"text":"([^"]+)"\}`)
	authorPattern := regexp.MustCompile(`"username":"([^"]+)"`)

	videos := []VideoMeta{}
	videoMatches := videoURLPattern.FindAllStringSubmatch(json, -1)
	thumbMatches := thumbPattern.FindAllStringSubmatch(json, -1)
	captionMatches := captionPattern.FindAllStringSubmatch(json, -1)
	authorMatches := authorPattern.FindAllStringSubmatch(json, -1)

	for i, match := range videoMatches {
		v := VideoMeta{
			VideoURL:  strings.ReplaceAll(match[1], "\\u0026", "&"),
			Thumbnail: safeGet(thumbMatches, i),
			Caption:   safeGet(captionMatches, i),
			Author:    safeGet(authorMatches, i),
		}
		videos = append(videos, v)
	}

	log.Printf("[parser] %d videos extracted\n", len(videos))
	return videos
}

func safeGet(matches [][]string, i int) string {
	if i < len(matches) && len(matches[i]) > 1 {
		return strings.ReplaceAll(matches[i][1], "\\u0026", "&")
	}
	return ""
}
