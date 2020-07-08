package main

import (
	"fmt"
	"sync"
)

// Fetcher represents an object that can fetch URLs' body.
type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

type crawled struct {
	urls map[string]bool
	mux  sync.Mutex
}

// Crawl uses fetcher to recursively crawl pages starting with url,
// to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	fetchedUrls := crawled{urls: make(map[string]bool)}
	var crawler func(string, int, Fetcher)
	crawler = func(url string, depth int, fetcher Fetcher) {
		var wg sync.WaitGroup

		fetchedUrls.mux.Lock()
		_, already := fetchedUrls.urls[url]
		if !already {
			fetchedUrls.urls[url] = true
		}
		fetchedUrls.mux.Unlock()

		if already || depth <= 0 {
			return
		}

		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("found: %s %q\n", url, body)
		for _, u := range urls {
			wg.Add(1)
			go func(u string) {
				defer wg.Done()
				crawler(u, depth-1, fetcher)
			}(u)
		}
		wg.Wait()
	}
	crawler(url, depth, fetcher)
}

func fetch() {
	Crawl("https://golang.org/", 4, fetcher)
}

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
