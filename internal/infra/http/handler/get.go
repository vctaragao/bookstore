package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/vctaragao/book-crud/internal/book"
)

type GetBook struct {
	bookService book.BookService
}

func NewGetBookHandler(bookService book.BookService) GetBook {
	return GetBook{
		bookService: bookService,
	}
}

func (h *GetBook) Handler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		log.Println("no book id sent")
		returnError(w, http.StatusBadRequest, "no book id sent")
		return
	}

	bookId, err := uuid.Parse(id)
	if err != nil {
		log.Println("unable to parse book id", err)
		returnError(w, http.StatusBadRequest, "unable to parse book id")
		return
	}

	book, err := h.bookService.GetBook(bookId)
	if err != nil {
		log.Printf("getting book: %v\n", err)

		if err == sql.ErrNoRows {
			returnError(w, http.StatusNotFound, "book not found")
			return
		}

		returnError(w, http.StatusInternalServerError, "unable to get book")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Printf("enconding book response: %v\n", err)
		returnError(w, http.StatusInternalServerError, "unable to enconde response")
		return
	}
}
