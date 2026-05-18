package tools

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"search_engine/internal/blobs"
	"search_engine/internal/indexers"
	"search_engine/internal/repository"
	"search_engine/internal/utils"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/gocolly/colly/v2"
)

type Crawler struct {
	c   *colly.Collector
	rep *repository.Repository
}

var wikipedia = indexers.GetWikipediaIndexer()

func InitCrawler(rep *repository.Repository) *Crawler {
	c := colly.NewCollector(
		colly.AllowedDomains(
			wikipedia.GetIndexer(),
		),
	)
	c.DisableCookies()
	c.AllowURLRevisit = false

	return &Crawler{c, rep}
}

// crawls into webpages, saves them internally and return the results
func (cr *Crawler) CrawlIntoIndexer(term string) (*blobs.BlobList, error) {
	searchUrl := wikipedia.Search(term)

	mdChan := make(chan *blobs.Blob, utils.MAX_CONCURRENT_REQUESTS)
	var wg sync.WaitGroup
	var atomicConcurrent atomic.Int32

	cr.c.OnHTML(".mw-content-ltr", func(h *colly.HTMLElement) {
		atomicConcurrent.Add(1)
		wg.Go(func() {
			if h.Response.StatusCode != http.StatusOK {
				return
			}

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

			b.Title = strings.ToLower(term)
			b.Datetime = time.Now().UTC()
			b.Folder = wikipedia.GetIndexer()
			b.URL = h.Request.URL.RawPath
			b.Body = markdown
			b.StemWords(string(markdown))

			// TODO: fix
			if selector := h.DOM.Find("meta[property=\"description\"]"); selector != nil {
				b.Description = selector.AttrOr("property", "Not found")
			}

			if selector := h.DOM.Find("meta[name='description']"); selector.Length() > 0 {
				b.Description = selector.AttrOr("content", "Not found")
			} else {
				b.Description = "Not found"
			}

			if err := cr.rep.SaveBlob(b, wikipedia, &markdown); err != nil {
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
				if atomicConcurrent.Load() >= utils.MAX_CONCURRENT_REQUESTS {
					return
				}
				h.Request.Visit(l)
			}
		})
	})

	maxQ := min(len(searchUrl), utils.MAX_CONCURRENT_REQUESTS)
	for _, s := range searchUrl[:maxQ] {
		if err := cr.c.Visit(s); err != nil {
			log.Printf("error while visiting %s: %s\n", searchUrl, err.Error())
		}
	}
	wg.Wait()
	close(mdChan)

	bl := blobs.BlobList{}
	for md := range mdChan {
		bl.AppendBlob(md)
	}

	return &bl, nil
}
