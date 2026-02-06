package models

import "gorm.io/gorm"

// Product 模型
type Product struct {
    gorm.Model
    Code  string `json:"code"`
    Price uint   `json:"price"`
}