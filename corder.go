package corder

import (
	"github.com/gocolly/colly/v2"
	"net/url"
	"sync"
	"time"
)

type Corder struct {
	requestCount int
	reqLock      sync.Mutex

	responseCount int
	resLock       sync.Mutex

	errors  map[error][]*url.URL
	errLock sync.Mutex

	recordTime time.Time
}

func NewCorder(c *colly.Collector) *Corder {
	corder := &Corder{
		requestCount:  0,
		responseCount: 0,
		errors:        make(map[error][]*url.URL),
	}
	c.OnRequest(func(request *colly.Request) {
		corder.reqLock.Lock()
		defer corder.reqLock.Unlock()
		corder.requestCount++
	})
	c.OnResponse(func(response *colly.Response) {
		corder.resLock.Lock()
		defer corder.resLock.Unlock()
		corder.responseCount++
	})
	c.OnError(func(response *colly.Response, err error) {
		corder.errLock.Lock()
		defer corder.errLock.Unlock()
		if corder.errors[err] == nil {
			corder.errors[err] = make([]*url.URL, 0)
		}
		corder.errors[err] = append(corder.errors[err], response.Request.URL)
	})
	corder.recordTime = time.Now()
	return corder
}

func (c *Corder) RecordTime() time.Time {
	return c.recordTime
}

func (c *Corder) RequestCount() int {
	return c.requestCount
}

func (c *Corder) ResponseCount() int {
	return c.responseCount
}

func (c *Corder) ErrorCount() int {
	errorCount := 0
	for _, urls := range c.errors {
		errorCount += len(urls)
	}
	return errorCount
}

func (c *Corder) Errors() map[error][]*url.URL {
	return c.errors
}

func (c *Corder) Reset() {
	c.reqLock.Lock()
	defer c.reqLock.Unlock()
	c.resLock.Lock()
	defer c.resLock.Unlock()
	c.errLock.Lock()
	defer c.errLock.Unlock()
	c.responseCount, c.responseCount = 0, 0
	c.errors = make(map[error][]*url.URL)
	c.recordTime = time.Now()
}
