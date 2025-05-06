package search

import (
	"context"
	"fmt"
)

// Service handles the business logic for search operations
type Service struct {
	repo Repository
}

// NewService creates a new search service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// SearchProducts searches for products with the given parameters
func (s *Service) SearchProducts(ctx context.Context, params SearchParams) (*SearchResponse, error) {
	// Validate parameters
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	// Perform search
	response, err := s.repo.Search(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("search products error: %w", err)
	}

	return response, nil
}

// GetFacets retrieves all available facets
func (s *Service) GetFacets(ctx context.Context) (map[string]interface{}, error) {
	facets, err := s.repo.GetFacets(ctx)
	if err != nil {
		return nil, fmt.Errorf("get facets error: %w", err)
	}

	return facets, nil
}

// IndexProduct indexes a single product
func (s *Service) IndexProduct(ctx context.Context, product interface{}) error {
	if err := s.repo.Index(ctx, product); err != nil {
		return fmt.Errorf("index product error: %w", err)
	}
	return nil
}

// BulkIndexProducts indexes multiple products
func (s *Service) BulkIndexProducts(ctx context.Context, products []interface{}) error {
	if err := s.repo.BulkIndex(ctx, products); err != nil {
		return fmt.Errorf("bulk index products error: %w", err)
	}
	return nil
}

// DeleteProduct deletes a product from the index
func (s *Service) DeleteProduct(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete product error: %w", err)
	}
	return nil
}

// UpdateProduct updates a product in the index
func (s *Service) UpdateProduct(ctx context.Context, id string, product interface{}) error {
	if err := s.repo.Update(ctx, id, product); err != nil {
		return fmt.Errorf("update product error: %w", err)
	}
	return nil
}
