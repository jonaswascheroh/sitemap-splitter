package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

// URLSet represents the top-level structure of a sitemap
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

// URL represents a single URL within a sitemap
type URL struct {
	Loc string `xml:"loc"`
}

func main() {
	sitemapURL := os.Args[1]
	sitemapContent := downloadSitemap(sitemapURL)
	sitemapName := extractSitemapName(sitemapURL)
	urlSet := parseSitemap(sitemapContent)
	splitAndSaveSitemaps(urlSet, 250, sitemapName) // Adjust 50 to your desired number of URLs per sitemap
}

func downloadSitemap(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func parseSitemap(content []byte) URLSet {
	var urlSet URLSet
	if err := xml.Unmarshal(content, &urlSet); err != nil {
		panic(err)
	}

	return urlSet
}

func extractSitemapName(sitemapURL string) string {
	parsedURL, err := url.Parse(sitemapURL)
	if err != nil {
		panic(err)
	}
	sitemapName := path.Base(parsedURL.Path)
	extension := path.Ext(sitemapName)
	if extension == "" {
		return sitemapName
	}
	return strings.TrimSuffix(sitemapName, extension)
}

func splitAndSaveSitemaps(urlSet URLSet, urlsPerSitemap int, sitemapName string) {
	totalURLs := len(urlSet.URLs)
	sitemapsNeeded := totalURLs / urlsPerSitemap
	if totalURLs%urlsPerSitemap != 0 {
		sitemapsNeeded++
	}
	err := os.MkdirAll("output", 0755)
	if err != nil {
		panic(err)
	}

	for i := 0; i < sitemapsNeeded; i++ {
		start := i * urlsPerSitemap
		end := start + urlsPerSitemap
		if end > totalURLs {
			end = totalURLs
		}

		subset := URLSet{URLs: urlSet.URLs[start:end]}
		content, err := xml.MarshalIndent(subset, "", "  ")
		if err != nil {
			panic(err)
		}

		fileName := "output/" + sitemapName + "_" + strconv.Itoa(i+1) + ".xml"
		if err := os.WriteFile(fileName, content, 0644); err != nil {
			panic(err)
		}
		fmt.Println("Saved:", fileName)
	}
}
