package corder

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"io"
	"net/url"
	"sync"
	"time"
)

type Corder struct {
	reqCount int
	reqLock  sync.Mutex

	resCount int
	resLock  sync.Mutex

	errs    map[string][]*url.URL
	errLock sync.Mutex

	startTime time.Time
}

func NewCorder(c *colly.Collector) *Corder {
	corder := &Corder{
		reqCount: 0,
		resCount: 0,
		errs:     make(map[string][]*url.URL),
	}
	c.OnRequest(func(request *colly.Request) {
		corder.reqLock.Lock()
		defer corder.reqLock.Unlock()
		corder.reqCount++
	})
	c.OnResponse(func(response *colly.Response) {
		corder.resLock.Lock()
		defer corder.resLock.Unlock()
		corder.resCount++
	})
	c.OnError(func(response *colly.Response, err error) {
		corder.errLock.Lock()
		defer corder.errLock.Unlock()
		if corder.errs[err.Error()] == nil {
			corder.errs[err.Error()] = make([]*url.URL, 0)
		}
		corder.errs[err.Error()] = append(corder.errs[err.Error()], response.Request.URL)
	})
	corder.startTime = time.Now()
	return corder
}

func (c *Corder) StartTime() time.Time {
	return c.startTime
}

func (c *Corder) RequestCount() int {
	return c.reqCount
}

func (c *Corder) ResponseCount() int {
	return c.resCount
}

func (c *Corder) ErrorCount() int {
	errorCount := 0
	for _, urls := range c.errs {
		errorCount += len(urls)
	}
	return errorCount
}

func (c *Corder) Errors() map[string][]*url.URL {
	return c.errs
}

func (c *Corder) Reset() {
	c.reqLock.Lock()
	defer c.reqLock.Unlock()
	c.resLock.Lock()
	defer c.resLock.Unlock()
	c.errLock.Lock()
	defer c.errLock.Unlock()
	c.resCount, c.resCount = 0, 0
	c.errs = make(map[string][]*url.URL)
	c.startTime = time.Now()
}

func (c *Corder) Print(writer io.Writer) {
	writer.Write([]byte("--------- Colly Corder ---------\n"))
	writer.Write([]byte(fmt.Sprintf("          cost: %v\n", time.Now().Sub(c.StartTime()))))
	writer.Write([]byte(fmt.Sprintf(" request count: %d\n", c.RequestCount())))
	writer.Write([]byte(fmt.Sprintf("response count: %d\n", c.ResponseCount())))
	writer.Write([]byte(fmt.Sprintf("   error count: %d\n", c.ErrorCount())))
	if c.ErrorCount() != 0 {
		writer.Write([]byte(fmt.Sprintf("\nerror detail:\n")))
		for err, urls := range c.Errors() {
			writer.Write([]byte(fmt.Sprintf("> [ %v ] (count: %d)\n", err, len(urls))))
			for _, u := range urls {
				writer.Write([]byte(fmt.Sprintf("  - %v\n", u)))
			}
		}
	}
	writer.Write([]byte("--------------------------------\n"))
}
