package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/vctaragao/book-crud/internal/book/entity"
)

type PostgresAdapter struct {
	db *sql.DB
}

func NewPostgresAdapter(db *sql.DB) PostgresAdapter {
	return PostgresAdapter{
		db: db,
	}
}

func (a *PostgresAdapter) CreateBook(book entity.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := a.db.ExecContext(
		ctx,
		"INSERT INTO books (id, title, price, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		book.ID.String(),
		book.Title,
		book.Price,
		book.CreatedAt,
		book.UpdatedAt,
	)

	return err
}

func (a *PostgresAdapter) GetBook(bookID uuid.UUID) (entity.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	book := entity.Book{}
	err := a.db.
		QueryRowContext(
			ctx,
			"SELECT id, title, price, created_at, updated_at FROM books WHERE id = $1",
			bookID,
		).Scan(&book.ID, &book.Title, &book.Price, &book.CreatedAt, &book.UpdatedAt)

	return book, err
}

func (a *PostgresAdapter) GetBooks(page, size int) ([]entity.Book, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	books := []entity.Book{}
	rows, err := a.db.QueryContext(ctx, "SELECT id, title, price FROM books LIMIT ? OFFSET ? ", size, page*size)

	for rows.Next() {
		book := entity.Book{}
		if err := rows.Scan(book.ID, book.Title, book.Price); err != nil {
			return []entity.Book{}, 0, err
		}

		books = append(books, book)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var total int
	err = a.db.QueryRowContext(ctx, "SELECT count(*) FROM books").Scan(&total)
	if err != nil {
		return []entity.Book{}, 0, err
	}

	return books, total, err
}
