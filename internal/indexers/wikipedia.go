package indexers

import (
	"fmt"
	"strings"

	"search_engine/internal/utils"
)

var INDEXER_WIKIPEDIA = types_of_indexers{
	Indexer:    "en.wikipedia.org",
	SearchTerm: "/wiki",
}

var api_rest_json string = "https://en.wikipedia.org/w/rest.php/v1/search/page"

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

func (i types_of_indexers) Search(term string) map[string]*IndexerExtraInfo {
	jsonRes := &wikipedia_json_response{}
	results := map[string]*IndexerExtraInfo{}

	err := utils.HttpGet(fmt.Sprintf("%s?q=%s&limit=%d", api_rest_json, term, utils.MAX_CONCURRENT_REQUESTS), jsonRes)
	if err != nil {
		return results
	}

	for _, url := range jsonRes.Pages {
		if url.Key == "" {
			continue
		}

		new_index_info := IndexerExtraInfo{
			URL: fmt.Sprintf("https://%s%s/%s", i.Indexer, i.SearchTerm, url.Key),
			Key: strings.ToLower(url.Key),
		}

		if url.Title != "" {
			new_index_info.Title = url.Title
		}

		if url.Description != "" {
			new_index_info.Description = url.Description
		}

		results[new_index_info.Key] = &new_index_info
	}

	return results
}

func (i types_of_indexers) GetAbsoluteIndexerURL() string {
	return fmt.Sprintf("https://%s", i.Indexer)
}
