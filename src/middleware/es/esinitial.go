package es

import (
	"log"

	"github.com/elastic/go-elasticsearch/v6"
)

// ES  client
var ES *elasticsearch.Client

func Connect(addr, esuser, espassword string) {
	addresses := []string{addr}
	config := elasticsearch.Config{
		Addresses: addresses,
		Username:  esuser,
		Password:  espassword,
		CloudID:   "",
		APIKey:    "",
	}
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	ES = es
}
