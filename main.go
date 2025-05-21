package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <website-url>")
		os.Exit(1)
	}

	baseURL := os.Args[1]
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
	}

	_, err := url.ParseRequestURI(baseURL)
	if err != nil {
		fmt.Printf("Invalid URL: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Checking links on: %s\n\n", baseURL)

	links, err := fetchLinks(baseURL)
	if err != nil {
		fmt.Printf("Error fetching links: %s\n", err)
		os.Exit(1)
	}

	for _, link := range links {
		status, err := checkLink(link, baseURL)
		if err != nil {
			fmt.Printf("❌ %s (Error: %s)\n", link, err)
			continue
		}

		if status == http.StatusOK {
			fmt.Printf("✅ %s\n", link)
		} else {
			fmt.Printf("❌ %s (Status: %d)\n", link, status)
		}
	}
}

func fetchLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	bodyStr := string(body)
	start := 0

	for {
		aStart := strings.Index(bodyStr[start:], "<a ")
		if aStart == -1 {
			break
		}
		aStart += start

		hrefStart := strings.Index(bodyStr[aStart:], "href=\"")
		if hrefStart == -1 {
			start = aStart + 1
			continue
		}
		hrefStart += aStart + 6

		hrefEnd := strings.Index(bodyStr[hrefStart:], "\"")
		if hrefEnd == -1 {
			start = hrefStart + 1
			continue
		}
		hrefEnd += hrefStart

		link := bodyStr[hrefStart:hrefEnd]
		if link != "" {
			links = append(links, link)
		}

		start = hrefEnd + 1
	}

	return links, nil
}

func checkLink(link, baseURL string) (int, error) {
	if strings.HasPrefix(link, "#") || 
	   strings.HasPrefix(link, "mailto:") || 
	   strings.HasPrefix(link, "tel:") ||
	   strings.HasPrefix(link, "javascript:") {
		return http.StatusOK, nil
	}

	u, err := url.Parse(link)
	if err != nil {
		return 0, err
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return 0, err
	}

	absURL := base.ResolveReference(u).String()

	if absURL == baseURL {
		return http.StatusOK, nil
	}

	resp, err := http.Head(absURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}