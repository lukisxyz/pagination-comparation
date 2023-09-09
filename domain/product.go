package domain

import (
	"time"
)

type Product struct {
	ID          int
	Name        string
	Amount      float64
	Description string
	CreatedAt   time.Time
}

type FetchProductsResponse[T int | string] struct {
	Error error
	Data  []Product
	Meta  Meta[T]
}
