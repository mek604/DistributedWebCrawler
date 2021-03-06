package main

import (
	"fmt"
	"net/url"
	"net/http"
	"golang.org/x/net/html"
	"strings"
)


func main() {

	// urls := [...] string {
	// 	"http://1867.com/",
	// 	"http://www.facebook.com/",
	// 	"https://www.facebook.com/",
	// 	"http://www.cs.ubc.ca/~bestchai/teaching/../teaching/././../teaching/cs416_2016w2/assign5/index.html",
	// 	"http://www.cs.ubc.ca/~bestchai/teaching/cs416_2016w2/assign5/index.html",
	// 	"https://www.49thapparel.com/",
	// }
	// fmt.Println("*** Get Domain Name Test *** ")
	// for _,u := range urls {
	// 	fmt.Println(getDomainName(u))
	// }

	// fmt.Println("\n*** Remove Relative Path Test ***")
	// for _,u := range urls {
	// 	fmt.Println(getAbsolutePath(u))

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

	// fmt.Println("*** Get DomainName Test ***")
	// for _,u := range urls {
	// 	fmt.Println(u, "\n\t", getDomainName(u))
	// }

	// fmt.Println("\n*** Remove Relative Path Test ***")
	// for _,u := range urls {
	// 	fmt.Println(resolveReference(u))
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
	
	// fmt.Println(crawl(urls[len(urls) - 1]))

	// fmt.Println("\n*** Formatting Address Test ***")
	// domain := getDomainName(urls[len(urls)-1])
	// links := []string {
	// 	"/collections/new-arrivals",
	// 	"/collections/./../collections/././new-arrivals",
	// 	"https://www.49thapparel.com/collections/././shoes",
	// 	"https://www.49thapparel.com/collections/boots",
	// 	"http://www.49thapparel.com/collections/boots",
	// 	"https://www.49thapparel.com/collections/./../collections/././new-arrivals",
	// }
	// for _,link := range links {
	// 	fmt.Println(filterAddress(domain, link))
	// }
	// domain := getDomainName(links[2])
	// webGraph := make(map[string][]string)
	// fmtRequestURL := "https://www.49thapparel.com/c/"
	// for _, link := range links {
	// 	if !contains(webGraph[fmtRequestURL], filterAddress(link, domain)) {
	// 		webGraph[fmtRequestURL] = append(webGraph[fmtRequestURL], filterAddress(link, domain))
	// 	}
	// }
	// for i,j := range webGraph {
	// 	fmt.Println(i, j)
	// }

	domain := getDomainName("https://www.facebook.com/")
	links := []string{
		"https://de-de.facebook.com/",
		"https://www.facebook.com/recover/initiate?lwv=110&ars=royal_blue_bar",
		"http://www.facebook.com/legal/terms/update",
		"www.facebook.com/legal/terms/update",
		"http://www.facebook.com/legal/terms/update",
		"/legal/terms/update",
		// "http://1867.com/",
		// "http://www.facebook.com/",
		// "https://www.facebook.com/",
		// "http://www.cs.ubc.ca/~bestchai/teaching/../teaching/././../teaching/cs416_2016w2/assign5/index.html",
		// "http://www.cs.ubc.ca/~bestchai/teaching/cs416_2016w2/assign5/index.html",
		// "https://www.49thapparel.com/",
		// "/collections/new-arrivals",
		// "/collections/./../collections/././new-arrivals",
	}
	for _, link := range links {
		x := filterAddress(link, domain) 
		fmt.Println(link, "\n\t", x)
		fmt.Println("\t\t", getDomainName(x))
	}


}

func getDomainName(uri string) string {
	u, _ := url.Parse(uri)
	return u.Host
}


func resolveReference(uri string) string {
	u, _ := url.Parse(uri)
	base, _ := url.Parse(getDomainName(uri))
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

// append to domain a relative path
// the input must be a domain (i.e, does not start with 'http')
func filterAddress(link, domain string) string {
	resolved := resolveReference(link)
	if strings.HasPrefix(link, "/") {
		resolved = "http://" + domain + resolved
		return resolved
	}
	return resolved
}
func contains(arr []string, val string) bool {
	for _, a := range arr {
		if val == a {
			return true
		}
	}
	return false
}