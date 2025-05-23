package api

import (
	"github.com/TheFellow/arch-lint/example/delta/bookstore/app/authors"
	"github.com/TheFellow/arch-lint/example/delta/bookstore/app/books"
)

type api struct {
}

func (a api) Book() books.Book {
	return books.Book{}
}

func (a api) Author() authors.Author {
	return authors.Author{}
}
