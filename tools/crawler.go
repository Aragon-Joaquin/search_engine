package tools

import (
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"search_engine/internal/blobs"
	"search_engine/internal/repository"
	"search_engine/internal/utils"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/gocolly/colly/v2"
)

const (
	MAX_CONCURRENT_REQUESTS = 1
)

type Crawler struct {
	c   *colly.Collector
	rep *repository.Repository
}

func InitCrawler(rep *repository.Repository) *Crawler {
	c := colly.NewCollector(
		colly.AllowedDomains(
			utils.INDEXER_WIKIPEDIA.Indexer,
		),
	)
	c.DisableCookies()
	c.AllowURLRevisit = false

	return &Crawler{c, rep}
}

// crawls into webpages, saves them internally and return the results
func (cr *Crawler) CrawlIntoIndexer(term string) []*blobs.Blob {
	searchUrl := utils.INDEXER_WIKIPEDIA.Search(term)

	mdChan := make(chan *blobs.Blob, MAX_CONCURRENT_REQUESTS)
	var wg sync.WaitGroup
	var atomicConcurrent atomic.Int32

	cr.c.OnHTML("body", func(h *colly.HTMLElement) {
		atomicConcurrent.Add(1)
		wg.Go(func() {
			// parse the content
			bodyNode := h.DOM.Nodes[h.Index]
			if bodyNode == nil {
				return
			}

			markdown, err := htmltomarkdown.ConvertNode(bodyNode)
			if err != nil {
				return
			}

			// send it to the channel
			b := blobs.CreateBlob()
			// pageTitle := h.ChildText(".mw-page-title-main") // wikipedia hardcodded

			b.Title = strings.ToLower(term)
			b.Datetime = time.Now().UTC()
			b.Folder = utils.INDEXER_WIKIPEDIA.Indexer
			b.URL = h.Request.URL.RawPath
			b.Body = markdown
			b.StemWords(string(markdown))

			// TODO: fix
			if selector := h.DOM.Find("meta[property=\"description\"]"); selector != nil {
				b.Description = selector.AttrOr("property", "Not found")
			}

			if err := cr.rep.SaveBlob(b, utils.INDEXER_WIKIPEDIA, &markdown); err != nil {
				return
			}

			mdChan <- b

			// look for more content
			var parseableLinks []string
			for _, url := range h.ChildAttrs("a[href]", "a") {
				res := h.Request.AbsoluteURL(url)

				if res == "" {
					continue
				}

				parseableLinks = append(parseableLinks, res)
			}

			// visit each one
			for _, l := range parseableLinks {
				if atomicConcurrent.Load() >= MAX_CONCURRENT_REQUESTS {
					return
				}
				h.Request.Visit(l)
			}
		})
	})

	var results []*blobs.Blob

	if err := cr.c.Visit(searchUrl); err != nil {
		log.Printf("error while visiting %s: %s\n", searchUrl, err.Error())
		return results
	}

	wg.Wait()
	close(mdChan)

	for md := range mdChan {
		results = append(results, md)
	}

	return results
}
