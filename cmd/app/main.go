package main

import (
	"fmt"
	"log"
	"os"

	models "github.com/pedronvasconcelos/golang-worker-crypto-price-historic/internal/model"
	"github.com/pedronvasconcelos/golang-worker-crypto-price-historic/internal/services"
)

func main() {

	resultsChan := make(chan models.CryptoPrice, len(workers))

	for _, worker := range workers {
		go func(w Worker) {
			price, err := executeWorker(w)
			if err == nil {
				resultsChan <- price
			}
		}(worker)
	}
	fmt.Println("Worker started")
	go func() {
		close(resultsChan)
	}()

	var coinPricesList []models.CryptoPrice
	for price := range resultsChan {
		coinPricesList = append(coinPricesList, price)
	}

	for _, price := range coinPricesList {
		fmt.Println(price)
	}
}
func executeWorker(worker Worker) (models.CryptoPrice, error) {
	switch worker.APISource {

	case "brapi":
		apiClient := services.NewBrapiClient(os.Getenv("BRAPI_API_KEY"))
		price, err := apiClient.GetBitcoinPrice()
		if err != nil {
			log.Fatalf("Failed to get bitcoin price: %v", err)
			return models.CryptoPrice{}, err
		}
		return price, nil

	}
	return models.CryptoPrice{}, fmt.Errorf("API source not found")
}

var workers = []Worker{
	{APISource: "brapi"},
	{APISource: "coinapi"},
}

type Worker struct {
	APISource string
}
