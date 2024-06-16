package book

import (
	"github.com/google/uuid"
	"github.com/vctaragao/book-crud/internal/book/entity"
)

type BookService struct {
	repo Repository
}

func NewBookService(repo Repository) BookService {
	return BookService{
		repo: repo,
	}
}

func (s *BookService) Create(title string, price float32) (entity.Book, error) {
	book := entity.NewBook(title, price)

	if err := book.Validate(); err != nil {
		return entity.Book{}, err
	}

	if err := s.repo.CreateBook(book); err != nil {
		return entity.Book{}, err
	}

	return book, nil
}

func (s *BookService) GetBook(bookId uuid.UUID) (entity.Book, error) {
	book, err := s.repo.GetBook(bookId)
	if err != nil {
		return entity.Book{}, err
	}

	return book, nil
}

func (s *BookService) GetBooks(page, size int) ([]entity.Book, int, error) {
	if page == 0 {
		page = 1
	}

	if size > 100 || size < 0 {
		size = 100
	}

	books, total, err := s.repo.GetBooks(page, size)
	if err != nil {
		return []entity.Book{}, 0, err
	}

	return books, total, nil
}
