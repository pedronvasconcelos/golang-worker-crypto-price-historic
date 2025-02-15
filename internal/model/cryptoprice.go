package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CryptoPrice struct {
	ID        uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	PriceUSD  decimal.Decimal `gorm:"type:numeric(20,8)" json:"priceUSD"`
	PriceTime time.Time       `gorm:"index" json:"priceTime,omitempty"`
	Coin      string          `gorm:"size:50;index" json:"coin"`
	Source    string          `gorm:"size:100" json:"source"`
}
