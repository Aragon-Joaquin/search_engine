package tools

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"search_engine/internal/blobs"
	"search_engine/internal/utils"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/gocolly/colly/v2"
)

const (
	MAX_CONCURRENT_REQUESTS = 10
)

// crawls into webpages, saves them internally and return the results
func CrawlIntoIndexer(term string) []*blobs.Blob {
	c := colly.NewCollector(
		colly.AllowedDomains(utils.GetAbsoluteIndexerURL(utils.INDEXER_WIKIPEDIA)),
	)
	c.DisableCookies()
	c.AllowURLRevisit = false

	searchUrl := utils.GetURL(utils.INDEXER_WIKIPEDIA, term)

	mdChan := make(chan *blobs.Blob, MAX_CONCURRENT_REQUESTS)
	var wg sync.WaitGroup
	var atomicConcurrent atomic.Int32

	c.OnHTML("body", func(h *colly.HTMLElement) {
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
			pageTitle := h.ChildText(".mw-page-title-main")

			b.Title = pageTitle
			b.Datetime = time.Now().UTC()
			b.Folder = utils.INDEXER_WIKIPEDIA

			if selector := h.DOM.Find("meta[property=\"description\"]"); selector != nil {
				b.Description = selector.AttrOr("property", "Not found")
			}

			mdChan <- b
			b.SaveBlob(utils.INDEXER_WIKIPEDIA, pageTitle, &markdown)

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

	if err := c.Visit(searchUrl); err != nil {
		log.Printf("error while visiting %s: %e\n", searchUrl, err)
		return results
	}

	wg.Wait()
	close(mdChan)

	for md := range mdChan {
		log.Println(md.Title)
		results = append(results, md)
	}

	return results
}
