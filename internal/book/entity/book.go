package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidTitle = errors.New("invalid book title")
var ErrInvalidPrice = errors.New("invalid book price")

type Book struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Price     float32   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBook(title string, price float32) Book {
	return Book{
		ID:        uuid.New(),
		Title:     title,
		Price:     price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (b *Book) Validate() error {
	if b.Title == "" {
		return ErrInvalidTitle
	}

	if b.Price < 0 {
		return ErrInvalidPrice
	}

	return nil
}
