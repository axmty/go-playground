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

// Waiter represents an object that can wait for goroutines to finish.
type Waiter struct {
	mux sync.Mutex
	n   int
	c   chan int
}

// NewWaiter creates a new Waiter object.
func NewWaiter() *Waiter {
	return &Waiter{n: 0, c: make(chan int)}
}

// Add adds n to the number of goroutines to wait for.
func (w *Waiter) Add(n int) {
	w.mux.Lock()
	defer w.mux.Unlock()
	w.n += n
}

// Wait waits for all goroutines to finish.
func (w *Waiter) Wait() {
	for range w.c {
		w.n--
		if w.n <= 0 {
			close(w.c)
		}
	}
}

// Done signals that a goroutine has finished.
func (w *Waiter) Done() {
	w.c <- 1
}

// Crawl uses fetcher to recursively crawl pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	fetchedUrls := crawled{urls: make(map[string]bool)}
	var crawler func(string, int, Fetcher)
	crawler = func(url string, depth int, fetcher Fetcher) {
		var w *Waiter = NewWaiter()

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
			w.Add(1)
			go func(u string) {
				defer w.Done()
				crawler(u, depth-1, fetcher)
			}(u)
		}
		w.Wait()
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
