package main

import (
	"fmt"
	"net/url"
	"net/http"
	"golang.org/x/net/html"
)


func main() {

	urls := [...] string {
		"http://1867.com/",
		"http://www.facebook.com/",
		"https://www.facebook.com/",
		"http://www.cs.ubc.ca/~bestchai/teaching/../teaching/././../teaching/cs416_2016w2/assign5/index.html",
		"http://www.cs.ubc.ca/~bestchai/teaching/cs416_2016w2/assign5/index.html",
		"https://www.49thapparel.com/",
	}
	// fmt.Println("*** Get Domain Name Test *** ")
	// for _,u := range urls {
	// 	fmt.Println(getDomainName(u))
	// }

	// fmt.Println("\n*** Remove Relative Path Test ***")
	// for _,u := range urls {
	// 	fmt.Println(getAbsolutePath(u))
	// }

	// mmap := make(map[string][]string) // maps each worker's ip:port to addresses
	// wip := "127.0.0.1:3000"
	// arr := []string{urls[0]}
	// mmap[wip] = []string{urls[1]}
	// if mmap[wip] == nil {
	// 	mmap[wip] = arr
	// } else {
	// 	mmap[wip] = append(mmap[wip], urls[1])
	// }
	// fmt.Println(mmap)


	// for _,u := range urls {
	// 	fmt.Println(getDomainName(u))
	// }

	// fmt.Println("\n*** Remove Relative Path Test ***")
	// for _,u := range urls {
	// 	fmt.Println(getAbsolutePath(u))
	// }

	// mmap := make(map[string][]string) // maps each worker's ip:port to addresses
	// wip := "127.0.0.1:3000"
	// arr := []string{urls[0]}
	// mmap[wip] = []string{urls[1]}
	// if mmap[wip] == nil {
	// 	mmap[wip] = arr
	// } else {
	// 	mmap[wip] = append(mmap[wip], urls[1])
	// }
	// fmt.Println(mmap)
	fmt.Println(crawl(urls[len(urls) - 1]))
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

// reference from extracting url links in golang
// [source] https://vorozhko.net/get-all-links-from-html-page-with-go-lang
func crawl(uri string) (links []string) {
	// depth := 1
	// depth loop here
	fmt.Println("Crawling", uri)
	resp, _ := http.Get(uri)
	defer resp.Body.Close()
	z := html.NewTokenizer(resp.Body)
	uniqueLinks := make(map[string]bool)

	for {
		tt := z.Next()
		switch tt {
			case html.ErrorToken:
				return
			case html.StartTagToken, html.EndTagToken:
				token := z.Token()
				if "a" == token.Data {
					// take add only unique links
					for _, attr := range token.Attr {
						if attr.Key == "href" && ! uniqueLinks[attr.Val] {
							links = append(links, attr.Val)
							uniqueLinks[attr.Val] = true
						}
					}
				}
		}
	}
	return
}
