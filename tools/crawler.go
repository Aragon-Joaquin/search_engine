package tools

import (
	"log"
	"net/http"
	"path"
	"strings"
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
		colly.Async(true),
	)
	c.DisableCookies()
	c.AllowURLRevisit = false

	return &Crawler{c, rep}
}

var searchUrl = map[string]*indexers.IndexerExtraInfo{}

// crawls into webpages, saves them internally and return the results
func (cr *Crawler) CrawlIntoIndexer(term string) (*blobs.BlobList, error) {
	searchUrl = wikipedia.Search(term)

	mdChan := make(chan *blobs.Blob, utils.MAX_CONCURRENT_REQUESTS)

	cr.c.OnHTML(".mw-content-ltr", func(h *colly.HTMLElement) {
		if h.Response.StatusCode != http.StatusOK {
			return
		}

		// match the blob -- needs to be improved though...
		path := path.Base(h.Request.URL.Path)
		info, ok := searchUrl[strings.ToLower(path)]
		if !ok {
			return
		}

		log.Println("BLOB SEARCH: ", info.Title)

		// parse the content
		if len(h.DOM.Nodes) == 0 {
			return
		}

		markdown, err := htmltomarkdown.ConvertNode(h.DOM.Get(0))
		if err != nil {
			return
		}

		// send it to the channel
		b := blobs.CreateBlob()

		b.Title = info.Title
		b.Datetime = time.Now().UTC()
		b.Folder = wikipedia.GetIndexer()
		b.URL = h.Request.URL.RawPath
		b.Body = markdown
		b.StemWords(string(markdown))

		// TODO: fix
		if info.Description == "" {
			if selector := h.DOM.Find("meta[property=\"description\"]"); selector != nil {
				b.Description = selector.AttrOr("property", "Not found")
			}

			if selector := h.DOM.Find("meta[name='description']"); selector.Length() > 0 {
				b.Description = selector.AttrOr("content", "Not found")
			} else {
				b.Description = "Not found"
			}
		} else {
			b.Description = info.Description
		}
		if err := cr.rep.SaveBlob(b, wikipedia, &markdown); err != nil {
			return
		}

		mdChan <- b
	})

	// WARN: THIS SHOULD NEVER EXCEED THE utils.MAX_CONCURRENT_REQUESTS. NEVER!!!
	for _, v := range searchUrl {

		// it should never happen. but whatever
		if v == nil {
			continue
		}

		if err := cr.c.Visit((*v).URL); err != nil {
			log.Printf("error while visiting %s: %s\n", (*v).URL, err.Error())
		}
	}

	cr.c.Wait()
	close(mdChan)

	bl := blobs.BlobList{}
	for md := range mdChan {
		bl.AppendBlob(md)
	}

	return &bl, nil
}
