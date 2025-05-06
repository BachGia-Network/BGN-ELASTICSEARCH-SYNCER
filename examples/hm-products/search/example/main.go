package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"your-module/search"

	"github.com/olivere/elastic/v7"
)

func main() {
	// Create Elasticsearch client
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9201"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("Error creating elasticsearch client: %v", err)
	}

	// Create repository
	repo := search.NewElasticsearchRepository(client, "products", 30)

	// Create service
	service := search.NewService(repo)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Example: Search products
	params := search.SearchParams{
		Query:    "T-Shirt",
		Page:     1,
		PageSize: 10,
		Filters: map[string]string{
			"material": "Cotton",
			"color":    "White",
		},
		SortField: "name",
		SortOrder: "asc",
	}

	results, err := service.SearchProducts(ctx, params)
	if err != nil {
		log.Fatalf("Error searching products: %v", err)
	}

	// Print results
	fmt.Printf("Found %d products\n", results.Total)
	for _, result := range results.Results {
		fmt.Printf("Product: %s (Score: %.2f)\n", result.Name, result.Score)
	}

	// Example: Get facets
	facets, err := service.GetFacets(ctx)
	if err != nil {
		log.Fatalf("Error getting facets: %v", err)
	}

	// Print facets
	fmt.Println("\nAvailable Facets:")
	for name, values := range facets {
		fmt.Printf("%s: %v\n", name, values)
	}
}
