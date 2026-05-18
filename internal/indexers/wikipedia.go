package indexers

import (
	"fmt"

	"search_engine/internal/utils"
)

var INDEXER_WIKIPEDIA = types_of_indexers{
	Indexer:    "en.wikipedia.org",
	SearchTerm: "/wiki",
}

var (
	api_rest_json       string  = "https://en.wikipedia.org/w/rest.php/v1/search/page"
	MIN_SCORE_THRESHOLD float64 = 0.05 // 5%
)

type wikipedia_json_response struct {
	Pages []page `json:"pages"`
}

type page struct {
	ID           int    `json:"id"`
	Key          string `json:"key"`
	Title        string `json:"title"`
	Excerpt      string `json:"excerpt"`
	MatchedTitle string `json:"matched_title"`
	Anchor       string `json:"anchor"`
	Description  string `json:"description"`
	Thumbnail    *struct {
		Mimetype string   `json:"mimetype"`
		Width    int      `json:"width"`
		Height   int      `json:"height"`
		Duration *float64 `json:"duration"`
		URL      string   `json:"url"`
	} `json:"thumbnail"`
}

func (i types_of_indexers) Search(term string) []string {
	var jsonRes *wikipedia_json_response
	results := []string{}

	err := utils.HttpGet(fmt.Sprintf("%s?q=%s&limit=%d", api_rest_json, term, 10), jsonRes)
	if err != nil || jsonRes == nil || jsonRes.Pages == nil {
		return results
	}

	for _, url := range jsonRes.Pages {
		if url.Key == "" {
			continue
		}
		results = append(results, url.Key)
	}

	return results
}

func (i types_of_indexers) GetAbsoluteIndexerURL() string {
	return fmt.Sprintf("https://%s", i.Indexer)
}
