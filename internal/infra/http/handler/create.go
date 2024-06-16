package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/vctaragao/book-crud/internal/book"
)

type CreateBook struct {
	bookService book.BookService
}

func NewCreateBookHandler(bookService book.BookService) CreateBook {
	return CreateBook{
		bookService: bookService,
	}
}

func (h *CreateBook) Handler(w http.ResponseWriter, r *http.Request) {
	type createBookParams struct {
		Title string  `json:"title"`
		Price float32 `json:"price"`
	}

	var params createBookParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Println("error on parsing body: ", err)

		if errors.Is(err, io.EOF) {
			returnError(w, http.StatusBadRequest, "empty body")
			return
		}

		returnError(w, http.StatusBadRequest, "unable to parse body")
		return
	}

	book, err := h.bookService.Create(params.Title, params.Price)
	if err != nil {
		log.Printf("creating book: %v\n", err)
		returnError(w, http.StatusInternalServerError, "unable to create book")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Printf("enconding book response: %v\n", err)
		returnError(w, http.StatusInternalServerError, "unable to enconde response")
		return
	}
}
