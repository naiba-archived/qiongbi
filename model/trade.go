package model

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

// Trade ..
type Trade struct {
	gorm.Model
	Name   string
	Email  string
	Note   string
	Amount decimal.Decimal `gorm:"type:decimal(20,2)"`
	Paid   bool
}
