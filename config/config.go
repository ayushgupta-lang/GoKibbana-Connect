package config

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/joho/godotenv"
)

type ElasticsearchConfig struct {
	ESClient *elasticsearch.Client
}

// NewElasticsearchClient initializes the Elasticsearch client
func NewElasticsearchClient() *elasticsearch.Client {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	esURL := os.Getenv("ELASTICSEARCH_URL")
	esUsername := os.Getenv("ELASTICSEARCH_USERNAME")
	esPassword := os.Getenv("ELASTICSEARCH_PASSWORD")

	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
		Username:  esUsername,
		Password:  esPassword,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %s", err)
	}
	defer res.Body.Close()

	log.Println("Connected to Elasticsearch")
	return es
}
