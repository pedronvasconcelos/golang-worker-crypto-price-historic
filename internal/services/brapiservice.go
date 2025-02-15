package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	models "github.com/pedronvasconcelos/golang-worker-crypto-price-historic/internal/model"
	"github.com/shopspring/decimal"
)

type BrapiClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewBrapiClient(token string) *BrapiClient {
	return &BrapiClient{
		BaseURL:    "https://api.brapi.dev/api/v2",
		APIKey:     token,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (bc *BrapiClient) GetBitcoinPrice() (models.CryptoPrice, error) {
	endpoint := "/crypto"
	url := fmt.Sprintf("%s%s?coin=BTC&currency=USD&token=%s", bc.BaseURL, endpoint, bc.APIKey)
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return models.CryptoPrice{}, fmt.Errorf("erro criando requisição: %v", err)
	}

	resp, err := bc.HTTPClient.Do(req)
	if err != nil {
		return models.CryptoPrice{}, fmt.Errorf("erro na requisição HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.CryptoPrice{}, fmt.Errorf("status code inválido: %d", resp.StatusCode)
	}

	var apiResponse BrapiCryptoResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return models.CryptoPrice{}, fmt.Errorf("erro decodificando JSON: %v", err)
	}

	if len(apiResponse.Coins) == 0 {
		return models.CryptoPrice{}, fmt.Errorf("nenhum dado de criptomoeda retornado")
	}

	btcData := apiResponse.Coins[0]
	priceTime := time.Unix(btcData.RegularMarketTime, 0)

	return models.CryptoPrice{
		ID:        uuid.New(),
		PriceUSD:  decimal.NewFromFloat(btcData.RegularMarketPrice),
		PriceTime: priceTime,
		Coin:      "BTC",
		Source:    "brapi.dev",
	}, nil
}

type BrapiCryptoResponse struct {
	Coins []BrapiCryptoData `json:"coins"`
}

type BrapiCryptoData struct {
	RegularMarketPrice float64 `json:"regularMarketPrice"`
	RegularMarketTime  int64   `json:"regularMarketTime"`
}
