package services

import (
	models "github.com/pedronvasconcelos/golang-worker-crypto-price-historic/internal/model"
)

type BitcoinService interface {
	GetBitcoinPrice() (models.CryptoPrice, error)
}
