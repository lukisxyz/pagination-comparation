package domain

type Meta[T int | string] struct {
	Total int `json:"total"`
	Limit int `json:"per_page"`
	Next  T   `json:"next"`
}
