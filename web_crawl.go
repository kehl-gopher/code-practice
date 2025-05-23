package main

import (
	"fmt"
	"sync"
)

/*
* In this exercise you'll use Go's concurrency features to parallelize a web crawler.
Modify the Crawl function to fetch URLs in parallel without fetching the same URL twice.
Hint: you can keep a cache of the URLs that have been fetched on a map, but maps alone are not safe for concurrent use!
*/
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

var urlLocked sync.RWMutex
var urlMap = make(map[string]bool)

func CheckUrlExists(url string) bool {
	urlLocked.Lock()
	defer urlLocked.Unlock()
	_, ok := urlMap[url]
	if !ok {
		urlMap[url] = true
		return ok
	}
	return ok
}

var mp sync.Map

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
// sync.Pool can be use to reuse object

func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}

	if _, ol := mp.LoadOrStore(url, true); ol {
		return
	}

	urlCh := make(chan map[string]interface{})
	errCh := make(chan error)

	go func() {
		var data map[string]interface{}
		body, urlss, err := fetcher.Fetch(url)
		if err != nil {
			errCh <- err
		}
		data = map[string]interface{}{
			"body": body,
			"urls": urlss,
		}
		urlCh <- data
	}()

	select {
	case err := <-errCh:
		fmt.Println(err)
	case uBod := <-urlCh:
		if uBod != nil {
		}
		fmt.Printf("found: %s %q\n", url, uBod["body"])
		urls := uBod["urls"].([]string)
		for _, u := range urls {
			Crawl(u, depth-1, fetcher)
		}
	}
}

// func main() {
// 	Crawl("https://golang.org/", 4, fetcher)
// }

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
