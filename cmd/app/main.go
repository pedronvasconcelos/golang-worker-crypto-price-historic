package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	models "github.com/pedronvasconcelos/golang-worker-crypto-price-historic/internal/model"
	"github.com/pedronvasconcelos/golang-worker-crypto-price-historic/internal/services"
)

var workers = []Worker{
	{APISource: "brapi"},
	{APISource: "coinapi"},
}

type Worker struct {
	APISource string
}

func main() {
	fmt.Println("Starting worker")
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Printf("Erro ao carregar o arquivo .env: %v\n", err)
		return // ou use os.Exit(1) se quiser encerrar o programa
	}

	resultsChan := make(chan models.CryptoPrice, len(workers))
	var wg sync.WaitGroup

	for _, worker := range workers {
		wg.Add(1)
		go func(w Worker) {
			defer wg.Done()
			price, err := executeWorker(w)
			if err == nil {
				resultsChan <- price
			}
			fmt.Printf("Worker finished: %v\n", w.APISource)
		}(worker)
	}

	go func() {
		wg.Wait()
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
	fmt.Printf("Executing worker: %v\n", worker.APISource)
	switch worker.APISource {
	case "brapi":
		apiClient := services.NewBrapiClient(os.Getenv("BRAPI_API_KEY"))
		fmt.Println("Getting bitcoin price from brapi")
		price, err := apiClient.GetBitcoinPrice()
		if err != nil {
			fmt.Printf("Failed to get bitcoin price: %v\n", err)
			return models.CryptoPrice{}, err
		}
		return price, nil
	}
	return models.CryptoPrice{}, fmt.Errorf("API source not found")
}
