package svc

import (
	"io"
	"net/http"

	"github.com/bacv/ethical_scraper/lib"
)

type Summary struct {
	lib.LinkSummary
	Success bool
	Error   string
}

type ToParse struct {
	Url  string
	Body string
}

type Scraper struct {
	httpClient   *http.Client
	httpQueue    chan ToParse
	httpPool     *lib.Pool[ToParse]
	parserPool   *lib.Pool[ToParse]
	summaryQueue chan Summary
}

func NewScraper(httpPoolSize, sumPoolSize uint32) *Scraper {
	httpClient := &http.Client{}
	s := &Scraper{
		httpClient: httpClient,
	}
	s.httpPool = lib.NewPool[ToParse](100, s.httpTask)
	s.parserPool = lib.NewPool[ToParse](2, s.parseTask)

	return s
}

func (s *Scraper) Start() <-chan Summary {
	s.httpPool.Start()
	s.parserPool.Start()
	return s.summaryQueue
}

func (s *Scraper) Scrape(url string) {
	s.httpPool.Do(ToParse{Url: url})
}

func (s *Scraper) Done() <-chan struct{} {
	return s.parserPool.Done()
}

func (s *Scraper) parseTask(p ToParse) {
	res := lib.CountLinks(p.Url, p.Body)

	sum := Summary{
		Success: true,
	}
	sum.LinkSummary = *res

	s.summaryQueue <- sum
}

func (s *Scraper) httpTask(p ToParse) {
	sum := Summary{}
	sum.PageUrl = p.Url

	resp, err := s.httpClient.Get(p.Url)
	if err != nil {
		sum.Error = err.Error()
		s.summaryQueue <- sum
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		sum.Error = "Invalid status code"
		s.summaryQueue <- sum
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		sum.Error = "Unable to parse response body"
		s.summaryQueue <- sum
		return
	}

	p.Body = string(body)
	s.parserPool.Do(p)
}
