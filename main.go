package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

// Function to scrape fonts from a given URL
func scrapeFonts(url string) {
	// Make HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find all <link> tags that refer to fonts (e.g., Google Fonts)
	doc.Find("link[rel='stylesheet']").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if exists && isFontURL(href) {
			fmt.Println("Font stylesheet URL found:", href)
		}
	})

	// Find all <style> tags and extract font URLs
	doc.Find("style").Each(func(index int, item *goquery.Selection) {
		css := item.Text()
		fontURLs := extractFontURLsFromCSS(css)
		for _, fontURL := range fontURLs {
			fmt.Println("Font URL found:", fontURL)
		}
	})
}

// Helper function to check if a URL is a font URL
func isFontURL(url string) bool {
	return strings.Contains(url, "fonts.googleapis.com") || strings.Contains(url, "font-awesome") || strings.Contains(url, "fonts.gstatic.com")
}

// Helper function to extract font URLs from CSS text
func extractFontURLsFromCSS(css string) []string {
	var fontURLs []string
	// Regular expression to find URLs within CSS content
	re := regexp.MustCompile(`url\(['"]?(https?://[^'"\)]+)['"]?\)`)
	matches := re.FindAllStringSubmatch(css, -1)
	for _, match := range matches {
		if len(match) > 1 {
			fontURLs = append(fontURLs, match[1])
		}
	}
	return fontURLs
}

func main() {
	url := "https://developer.fedoraproject.org/tech/languages/go/go-programs.html" // Replace with the URL you want to scrape
	scrapeFonts(url)
}
