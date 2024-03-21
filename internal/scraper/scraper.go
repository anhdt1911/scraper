package scraper

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	googleDomain = "https://www.google.com/search?q="
)

var userAgents = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
}

type Scraper struct {
	client *http.Client
}

type SearchResult struct {
	keyword           string
	links             []string
	totalSearchResult string
	htmlContent       string
	adwords           []string
}

func New() *Scraper {
	return &Scraper{
		http.DefaultClient,
	}
}

func buildUrl(keyword string) string {
	keyword = strings.TrimSpace(keyword)
	rgx := regexp.MustCompile(`\s+`)
	keyword = rgx.ReplaceAllString(keyword, "+")
	return fmt.Sprintf("%s%s", googleDomain, keyword)
}

func rotateUserAgent() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	randNum := rng.Int() % len(userAgents)
	return userAgents[randNum]
}

func (s *Scraper) request(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return nil, err
	}
	req.Header.Set("User-Agent", rotateUserAgent())

	res, err := s.client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return nil, err
	}
	return res, nil
}

func (s *Scraper) Scrape(keyword string) error {
	res, err := s.request(buildUrl(keyword))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		return err
	}

	// Populate result.
	var links []string
	// var totalSearchResult string
	var link func(*html.Node)
	link = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				for _, a := range n.Attr {
					// Filter out valid links.
					if a.Key == "href" && strings.HasPrefix(a.Val, "http") {
						links = append(links, a.Val)
					}
				}
			case "div":
				for _, a := range n.Attr {
					if strings.Contains(a.Val, "result-stats") {
						fmt.Println(a)
					}
				}
			}
		}

		if n.Type == html.TextNode {
			fmt.Println(n.Data)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			link(c)
		}
	}
	link(doc)

	return nil
}
