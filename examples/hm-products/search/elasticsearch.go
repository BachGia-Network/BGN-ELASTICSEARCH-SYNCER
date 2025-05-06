package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/olivere/elastic/v7"
)

// ElasticsearchRepository implements the Repository interface using Elasticsearch
type ElasticsearchRepository struct {
	client  *elastic.Client
	index   string
	timeout int
}

// NewElasticsearchRepository creates a new Elasticsearch repository
func NewElasticsearchRepository(client *elastic.Client, index string, timeout int) *ElasticsearchRepository {
	return &ElasticsearchRepository{
		client:  client,
		index:   index,
		timeout: timeout,
	}
}

// Search implements the Searcher interface
func (r *ElasticsearchRepository) Search(ctx context.Context, params SearchParams) (*SearchResponse, error) {
	// Create search service
	searchService := r.client.Search().
		Index(r.index).
		TimeoutInSeconds(r.timeout)

	// Build query
	query := elastic.NewBoolQuery()

	// Add text search if query is provided
	if params.Query != "" {
		query = query.Must(
			elastic.NewMultiMatchQuery(params.Query, "name^3", "categories.name^2", "attributes.value"),
		)
	}

	// Add filters
	for field, value := range params.Filters {
		switch field {
		case "category":
			query = query.Filter(elastic.NewTermQuery("categories.name.keyword", value))
		case "material":
			query = query.Filter(elastic.NewNestedQuery("attributes",
				elastic.NewBoolQuery().
					Must(
						elastic.NewTermQuery("attributes.name.keyword", "Material"),
						elastic.NewTermQuery("attributes.value.keyword", value),
					),
			))
		case "color":
			query = query.Filter(elastic.NewNestedQuery("attributes",
				elastic.NewBoolQuery().
					Must(
						elastic.NewTermQuery("attributes.name.keyword", "Color"),
						elastic.NewTermQuery("attributes.value.keyword", value),
					),
			))
		}
	}

	// Add aggregations
	searchService = searchService.Aggregation("categories", elastic.NewTermsAggregation().
		Field("categories.name.keyword").
		Size(20))

	searchService = searchService.Aggregation("materials", elastic.NewNestedAggregation().
		Path("attributes").
		SubAggregation("filter", elastic.NewFilterAggregation().
			Filter(elastic.NewTermQuery("attributes.name.keyword", "Material")).
			SubAggregation("values", elastic.NewTermsAggregation().
				Field("attributes.value.keyword").
				Size(20))))

	searchService = searchService.Aggregation("colors", elastic.NewNestedAggregation().
		Path("attributes").
		SubAggregation("filter", elastic.NewFilterAggregation().
			Filter(elastic.NewTermQuery("attributes.name.keyword", "Color")).
			SubAggregation("values", elastic.NewTermsAggregation().
				Field("attributes.value.keyword").
				Size(20))))

	// Add sorting
	if params.SortField != "" {
		order := elastic.Asc
		if strings.ToLower(params.SortOrder) == "desc" {
			order = elastic.Desc
		}
		searchService = searchService.Sort(params.SortField, order)
	}

	// Add pagination
	from := (params.Page - 1) * params.PageSize
	searchService = searchService.From(from).Size(params.PageSize)

	// Execute search
	result, err := searchService.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch search error: %w", err)
	}

	// Process results
	response := &SearchResponse{
		Total:        result.TotalHits(),
		Page:         params.Page,
		PageSize:     params.PageSize,
		Results:      make([]SearchResult, 0),
		Aggregations: result.Aggregations,
	}

	// Parse hits
	for _, hit := range result.Hits.Hits {
		var result SearchResult
		if err := json.Unmarshal(hit.Source, &result); err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}
		result.Score = *hit.Score
		result.RawSource = hit.Source
		response.Results = append(response.Results, result)
	}

	return response, nil
}

// GetFacets implements the Searcher interface
func (r *ElasticsearchRepository) GetFacets(ctx context.Context) (map[string]interface{}, error) {
	searchService := r.client.Search().
		Index(r.index).
		Size(0).
		TimeoutInSeconds(r.timeout)

	// Add all aggregations
	searchService = searchService.Aggregation("categories", elastic.NewTermsAggregation().
		Field("categories.name.keyword").
		Size(20))

	searchService = searchService.Aggregation("materials", elastic.NewNestedAggregation().
		Path("attributes").
		SubAggregation("filter", elastic.NewFilterAggregation().
			Filter(elastic.NewTermQuery("attributes.name.keyword", "Material")).
			SubAggregation("values", elastic.NewTermsAggregation().
				Field("attributes.value.keyword").
				Size(20))))

	searchService = searchService.Aggregation("colors", elastic.NewNestedAggregation().
		Path("attributes").
		SubAggregation("filter", elastic.NewFilterAggregation().
			Filter(elastic.NewTermQuery("attributes.name.keyword", "Color")).
			SubAggregation("values", elastic.NewTermsAggregation().
				Field("attributes.value.keyword").
				Size(20))))

	result, err := searchService.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch facets error: %w", err)
	}

	return result.Aggregations, nil
}

// Index implements the Indexer interface
func (r *ElasticsearchRepository) Index(ctx context.Context, document interface{}) error {
	_, err := r.client.Index().
		Index(r.index).
		BodyJson(document).
		Do(ctx)
	return err
}

// BulkIndex implements the Indexer interface
func (r *ElasticsearchRepository) BulkIndex(ctx context.Context, documents []interface{}) error {
	bulk := r.client.Bulk()
	for _, doc := range documents {
		req := elastic.NewBulkIndexRequest().
			Index(r.index).
			Doc(doc)
		bulk.Add(req)
	}
	_, err := bulk.Do(ctx)
	return err
}

// Delete implements the Indexer interface
func (r *ElasticsearchRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Delete().
		Index(r.index).
		Id(id).
		Do(ctx)
	return err
}

// Update implements the Indexer interface
func (r *ElasticsearchRepository) Update(ctx context.Context, id string, document interface{}) error {
	_, err := r.client.Update().
		Index(r.index).
		Id(id).
		Doc(document).
		Do(ctx)
	return err
}
