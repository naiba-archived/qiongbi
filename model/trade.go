package model

import "github.com/jinzhu/gorm"

// Trade ..
type Trade struct {
	gorm.Model
	Name   string
	Email  string
	Note   string
	Amount string
	Paid   bool
}
