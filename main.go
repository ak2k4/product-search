package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Product defines a product with searchable fields
type Product struct {
	ID       int
	Name     string
	Category string
}

var (
	products []Product   // In-memory product data
	index    bleve.Index // Bleve search index
)

func main() {
	// Step 1: Generate dummy data
	categories := []string{"Electronics", "Books", "Clothing", "Toys"}
	products = make([]Product, 1000000)
	for i := range products {
		products[i] = Product{
			ID:       i + 1,
			Name:     fmt.Sprintf("Product %d", i+1),
			Category: categories[i%len(categories)],
		}
	}

	// Step 2: Open or create Bleve index
	var err error
	index, err = bleve.Open("products.bleve")
	if err == bleve.ErrorIndexPathDoesNotExist {
		fmt.Println("Creating new Bleve index...")
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New("products.bleve", mapping)
		if err != nil {
			log.Fatal(err)
		}

		// Batch indexing for performance
		fmt.Println("Indexing products... (this may take a few minutes)")
		batch := index.NewBatch()
		for i, p := range products {
			id := strconv.Itoa(p.ID)
			batch.Index(id, p)
			if i > 0 && i%10000 == 0 {
				if err := index.Batch(batch); err != nil {
					log.Fatal(err)
				}
				batch = index.NewBatch()
			}
		}
		if err := index.Batch(batch); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Indexing complete.")
	} else if err != nil {
		log.Fatal(err)
	}

	// Step 3: Setup Chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5)) // Enable gzip compression

	// Health route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	// Search route
	r.Get("/search", searchHandler)

	// Step 4: Setup graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("Server started at http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	<-stop
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if err := index.Close(); err != nil {
		log.Printf("Error closing index: %v", err)
	}

	fmt.Println("Server stopped gracefully.")
}

// searchHandler handles GET /search?q=term[&from=0&size=50]
func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "Missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	// Support pagination
	fromStr := r.URL.Query().Get("from")
	sizeStr := r.URL.Query().Get("size")
	from, _ := strconv.Atoi(fromStr)
	size, _ := strconv.Atoi(sizeStr)
	if size == 0 {
		size = 50
	}

	// Use query string query for fielded search
	query := bleve.NewQueryStringQuery(q)
	searchRequest := bleve.NewSearchRequestOptions(query, size, from, false)

	searchResult, err := index.Search(searchRequest)
	if err != nil {
		http.Error(w, "Search error", http.StatusInternalServerError)
		return
	}

	matchedProducts := []Product{}
	for _, hit := range searchResult.Hits {
		id, err := strconv.Atoi(hit.ID)
		if err == nil && id > 0 && id <= len(products) {
			matchedProducts = append(matchedProducts, products[id-1])
		}
	}

	w.Header().Set("Content-Type", "application/json")
	prettyJSON, err := json.MarshalIndent(matchedProducts, "", "  ")
	if err != nil {
		http.Error(w, "Error formatting JSON", http.StatusInternalServerError)
		return
	}
	w.Write(prettyJSON)
}
