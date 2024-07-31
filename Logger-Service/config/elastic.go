package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"logger/models"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/joho/godotenv"
)

var Client *elasticsearch.Client

// ConnectElastic initializes the Elasticsearch client
func ConnectElastic() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env dosyası yüklenemedi:", err)
		panic(err)
	}

	elasticURL := os.Getenv("ELASTIC")
	if elasticURL == "" {
		log.Fatal("ELASTIC ortam değişkeni ayarlanmamış")
	}

	Client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{elasticURL},
	})
	if err != nil {
		log.Fatalf("Elasticsearch istemcisi oluşturulamadı: %v", err)
	}
	log.Println("Elasticsearch istemcisi oluşturuldu:", Client)
}

// SendLogToElasticsearch sends a log message to Elasticsearch
func SendLogToElasticsearch(logMessage models.LogMessage) error {
	log.Println("rabbitten gelen vahiy")
	if Client == nil {
		return fmt.Errorf("elasticsearch istemcisi oluşturulmamış")
	}
	log.Println(logMessage)
	logEl, err := json.Marshal(logMessage)
	if err != nil {
		log.Println("Log mesajı JSON'a dönüştürülemedi:", err)
		return err
	}

	req := esapi.IndexRequest{
		Index:   "logs",
		Body:    strings.NewReader(string(logEl)),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), Client)
	if err != nil {
		log.Println("Elasticsearch'e istek gönderilemedi:", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Println("Elasticsearch hata yanıtı:", res.String())
		return fmt.Errorf("elasticsearch hata yaniti: %s", res.Status())
	}
	log.Println("elasticsearch yanıtı:", res)
	return nil
}
