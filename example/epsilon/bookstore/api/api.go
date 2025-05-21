package api

import (
	"github.com/TheFellow/go-arch-lint/example/beta/bookstore/app/authors"
	"github.com/TheFellow/go-arch-lint/example/beta/bookstore/app/books"
)

type api struct {
}

func (a api) Book() books.Book {
	return books.Book{}
}

func (a api) Author() authors.Author {
	return authors.Author{}
}
