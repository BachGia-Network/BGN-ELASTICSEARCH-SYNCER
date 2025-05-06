package search

import (
	"context"
)

// SearchResult represents a single search result
type SearchResult struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Categories []Category             `json:"categories"`
	Variants   []Variant              `json:"variants"`
	Attributes []Attribute            `json:"attributes"`
	Score      float64                `json:"_score"`
	RawSource  map[string]interface{} `json:"-"`
}

// Category represents a product category
type Category struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
}

// Variant represents a product variant
type Variant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Attribute represents a product attribute
type Attribute struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// SearchParams represents search parameters
type SearchParams struct {
	Query     string            `json:"query"`
	Filters   map[string]string `json:"filters"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	SortField string            `json:"sort_field"`
	SortOrder string            `json:"sort_order"`
}

// SearchResponse represents the search response
type SearchResponse struct {
	Results      []SearchResult `json:"results"`
	Total        int64          `json:"total"`
	Page         int            `json:"page"`
	PageSize     int            `json:"page_size"`
	Aggregations interface{}    `json:"aggregations"`
}

// Searcher defines the interface for search operations
type Searcher interface {
	Search(ctx context.Context, params SearchParams) (*SearchResponse, error)
	GetFacets(ctx context.Context) (map[string]interface{}, error)
}

// Indexer defines the interface for index operations
type Indexer interface {
	Index(ctx context.Context, document interface{}) error
	BulkIndex(ctx context.Context, documents []interface{}) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, document interface{}) error
}

// Repository combines search and index operations
type Repository interface {
	Searcher
	Indexer
}
