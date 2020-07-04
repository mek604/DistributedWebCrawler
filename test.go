package main

import (
	"fmt"
	"net/url"
)


func main() {

	urls := [...] string {
		"http://1867.com/",
		"http://www.facebook.com/",
		"https://www.facebook.com/",
		"http://www.cs.ubc.ca/~bestchai/teaching/../teaching/././../teaching/cs416_2016w2/assign5/index.html",
		"http://www.cs.ubc.ca/~bestchai/teaching/cs416_2016w2/assign5/index.html",
	}

	fmt.Println("*** Get Domain Name Test *** ")
	for _,u := range urls {
		fmt.Println(getDomainName(u))
	}

	fmt.Println("\n*** Remove Relative Path Test ***")
	for _,u := range urls {
		fmt.Println(getAbsolutePath(u))
	}

	mmap := make(map[string][]string) // maps each worker's ip:port to addresses
	wip := "127.0.0.1:3000"
	arr := []string{urls[0]}
	mmap[wip] = []string{urls[1]}
	if mmap[wip] == nil {
		mmap[wip] = arr
	} else {
		mmap[wip] = append(mmap[wip], urls[1])
	}
	fmt.Println(mmap)
}

func getDomainName(uri string) string {
	u, _ := url.Parse(uri)
	return u.Host
}

func getAbsolutePath(uri string) string {
	u, _ := url.Parse(uri)
	base, _ := url.Parse(getDomainName("http://example.com/directory/"))
	return base.ResolveReference(u).String()
}