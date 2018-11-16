package spider

import (
	"fmt"
	"path"

	//"strings"

	"github.com/gocolly/colly"
)

var (
	defaultUrl = "http://172.16.33.11/Person/"
	defaultMaxDepth = 5
	defaultMatchSelector = "a[href$='.dcm']"
	defaultConcurrence = 2
	defaultOutput = "./output"
)

func InitSpider(url, output string, maxDepth, concurrence int) {
	if url == "" {
		url = defaultUrl
	}
	if output == "" {
		output = defaultOutput
	}
	if maxDepth == 0 {
		maxDepth = defaultMaxDepth
	}
	if concurrence == 0 {
		concurrence = defaultConcurrence
	}
	// Instantiate default collector
	c := colly.NewCollector(
		// MaxDepth is 1, so only the links on the scraped page
		// is visited, and no further links are followed
		colly.MaxDepth(maxDepth),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: concurrence})
	// On every a element which has href attribute call callback
	c.OnHTML("a[href*='.13777454333330444454.']", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnHTML("a[href*='I/']", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnHTML(defaultMatchSelector, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link


		//e.Request.Visit(link)
		if isFileExist(path.Join(output, "/data"), link) {
			return
		}
		fmt.Printf("Link found: -> %s\n", link)
		absoluteURL := e.Request.AbsoluteURL(link)
		//err := appendUrlToLocal(output, "url.txt", []byte(fmt.Sprintf("%s\r\n", absoluteURL)))
		//if err != nil {
		//	fmt.Println(link, "===>", err.Error())
		//	return
		//}
		data, err := download(absoluteURL)
		if err != nil {
			fmt.Println(link, "===>", err.Error())
			return
		}
		err = saveToLocal(path.Join(output, "/data"), link, data)
		if err != nil {
			fmt.Println(link, "===>", err.Error())
			return
		}
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(url)
	c.Wait()
}