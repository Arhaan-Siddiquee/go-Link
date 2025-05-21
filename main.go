package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
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

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links, nil
}

func checkLink(link, baseURL string) (int, error) {
	// Skip mailto: and other non-http links
	if strings.HasPrefix(link, "#") || strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "tel:") {
		return http.StatusOK, nil
	}

	// Handle relative URLs
	if strings.HasPrefix(link, "/") {
		parsedBase, err := url.Parse(baseURL)
		if err != nil {
			return 0, err
		}
	link = parsedBase.Scheme + "://" + parsedBase.Host + link
	} else if !strings.HasPrefix(link, "http") {
		// Handle relative paths without leading slash
		parsedBase, err := url.Parse(baseURL)
		if err != nil {
			return 0, err
		}
		link = parsedBase.Scheme + "://" + parsedBase.Host + "/" + link
	}

	// Skip checking the same base URL to avoid infinite recursion
	if link == baseURL {
		return http.StatusOK, nil
	}

	resp, err := http.Head(link) // Using HEAD to be more efficient than GET
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}