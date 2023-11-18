package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type ConcurrentMap struct {
	m     map[string]struct{}
	mutex sync.Mutex
}

func (m *ConcurrentMap) Add(s string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, ok := m.m[s]
	if ok {
		return false
	} else {
		m.m[s] = struct{}{}
		return true
	}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, cm *ConcurrentMap, ret chan string) {
	defer close(ret)
	if depth <= 0 {
		return
	}

	if !cm.Add(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		ret <- err.Error()
		return
	}

	ret <- fmt.Sprintf("found: %s %q\n", url, body)

	result := make([]chan string, len(urls))
	for i, u := range urls {
		result[i] = make(chan string)
		go Crawl(u, depth-1, fetcher, cm, result[i])
	}

	for i := range result {
		for s := range result[i] {
			ret <- s
		}
	}
}

func main() {
	result := make(chan string)
	cm := ConcurrentMap{m: make(map[string]struct{}), mutex: sync.Mutex{}}
	go Crawl("https://golang.org/", 4, fetcher, &cm, result)

	for s := range result {
		fmt.Println(s)
	}
	fmt.Println(cm)
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
