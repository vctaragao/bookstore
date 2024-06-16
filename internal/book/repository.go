package book

import (
	"github.com/google/uuid"
	"github.com/vctaragao/book-crud/internal/book/entity"
)

type (
	Repository interface {
		CreateBook(book entity.Book) error
		GetBook(bookID uuid.UUID) (entity.Book, error)
		GetBooks(page, size int) ([]entity.Book, int, error)
	}
)
