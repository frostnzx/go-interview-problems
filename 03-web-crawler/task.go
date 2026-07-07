package main

import (
	"fmt"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, ch chan map[string]bool, flg chan bool) {
	defer func() {
		flg <- true
	}()

	if depth <= 0 {
		return
	}

	// check visited
	// note : can use mutex instead
	vis := <-ch
	if vis[url] {
		ch <- vis
		return
	}
	vis[url] = true
	ch <- vis

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)

	flagWaitList := make([]chan bool, 0)
	for _, u := range urls {
		childFlg := make(chan bool, 1)
		flagWaitList = append(flagWaitList, childFlg)
		go Crawl(u, depth-1, fetcher, ch, childFlg)
	}
	for _, childFlg := range flagWaitList { // can use waitgroup instead
		<-childFlg // wait for all child to finish
	}
}
