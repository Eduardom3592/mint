package api

import (
	"fmt"
	"net/url"

	"github.com/programmersd21/mint/internal/models"
)

type SearchSort string

const (
	SortRelevance SearchSort = "relevance"
	SortDownloads SearchSort = "downloads"
	SortFollowers SearchSort = "followers"
	SortUpdated   SearchSort = "updated"
	SortCreated   SearchSort = "created"
	SortName      SearchSort = "title"
)

type SearchFilter struct {
	Query        string
	ProjectType  string
	Loaders      []string
	GameVersions []string
	Categories   []string
	Sort         SearchSort
	Offset       int
	Limit        int
	Facets       bool
}

func (c *Client) Search(filter SearchFilter) (*models.SearchResponse, error) {
	query := url.Values{}

	if filter.Query != "" {
		query.Set("query", filter.Query)
	}

	var facets [][]string

	if filter.ProjectType != "" {
		facets = append(facets, []string{fmt.Sprintf("project_type:%s", filter.ProjectType)})
	}

	for _, loader := range filter.Loaders {
		facets = append(facets, []string{fmt.Sprintf("categories:%s", loader)})
	}

	for _, version := range filter.GameVersions {
		facets = append(facets, []string{fmt.Sprintf("versions:%s", version)})
	}

	for _, cat := range filter.Categories {
		facets = append(facets, []string{fmt.Sprintf("categories:%s", cat)})
	}

	if len(facets) > 0 {
		facetJSON := "["
		for i, or := range facets {
			if i > 0 {
				facetJSON += ","
			}
			facetJSON += "[\"" + or[0] + "\"]"
		}
		facetJSON += "]"
		query.Set("facets", facetJSON)
	}

	if filter.Offset > 0 {
		query.Set("offset", fmt.Sprintf("%d", filter.Offset))
	}
	if filter.Limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", filter.Limit))
	} else {
		query.Set("limit", "50")
	}

	query.Set("index", string(filter.Sort))

	var result models.SearchResponse
	if err := c.get("/search", query, &result); err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	return &result, nil
}
