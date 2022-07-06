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

type VisitedUrl struct {
	mu   sync.Mutex
	urls []string
}

func (v *VisitedUrl) isExisted(url string) bool {
	for _, v := range v.urls {
		if v == url {
			return true
		}
	}
	return false
}

func (v *VisitedUrl) add(url string) {
	v.mu.Lock()
	v.urls = append(v.urls, url)
	v.mu.Unlock()
}

func (v *VisitedUrl) delete(url string) (string, bool) {
	v.mu.Lock()
	index := -1
	for i, u := range v.urls {
		if u == url {
			index = i
			break
		}
	}
	if index == -1 {
		return "", false
	}
	//swap element to last
	v.urls[index] = v.urls[len(v.urls)-1]
	//make new slice exclude last element
	v.urls = v.urls[:len(v.urls)-1]

	defer v.mu.Unlock()
	return url, true
}

var visitedUrl VisitedUrl

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	defer wait.Done()
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	if visitedUrl.isExisted(url) {
		return
	}
	// This implementation doesn't do either:
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	visitedUrl.add(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		wait.Add(1)
		go Crawl(u, depth-1, fetcher)
	}
	return
}

var wait sync.WaitGroup

func main() {
	wait.Add(1)
	go Crawl("https://golang.org/", 4, fetcher)
	wait.Wait()
	// var v VisitedUrl;
	// fmt.Println(v.isExisted("a"))
	// v.urls = make([]string,10)
	// for i := 0; i < 9; i++ {
	// 	v.add(fmt.Sprint(i))
	// }
	// fmt.Println(v.urls)
	// fmt.Printf("Info: %d, %p \n",len(v.urls),v.urls)
	// v.delete("3")
	// fmt.Println(v.urls)
	// fmt.Printf("Info: %d, %p \n",len(v.urls),v.urls)

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
