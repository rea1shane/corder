package corder

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"testing"
	"time"
)

func TestNewCorder(t *testing.T) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.baidu.com"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Response %s: %d bytes\n", r.Request.URL, len(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error %s: %v\n", r.Request.URL, err)
	})

	corder := NewCorder(c)

	c.Visit("http://www.baidu.com/")

	fmt.Println(" request count:", corder.RequestCount())
	fmt.Println("response count:", corder.ResponseCount())
	fmt.Println("   error count:", corder.ErrorCount())
	fmt.Println("          cost:", time.Now().Sub(corder.RecordTime()))
	fmt.Println(corder.Errors())
}
