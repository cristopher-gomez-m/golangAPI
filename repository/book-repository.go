package repository

import (
	"github.com/cristopher-gomez-m/golang_api/entity"
	"gorm.io/gorm"
)

type BookRepository interface {
	InsertBook(book entity.Book) entity.Book
	UpdateBook(book entity.Book) entity.Book
	DeleteBook(book entity.Book)
	AllBook() []entity.Book
	FindBookByID(bookId uint64) entity.Book
}

type bookConnection struct {
	connection *gorm.DB
}

func NewBookConnection(db *gorm.DB) BookRepository {
	return &bookConnection{
		connection: db,
	}
}

func (db *bookConnection) InsertBook(book entity.Book) entity.Book {
	db.connection.Save(&book)
	db.connection.Preload("User").Find(&book)
	return book
}
func (db *bookConnection) UpdateBook(book entity.Book) entity.Book {
	return db.InsertBook(book)
}
func (db *bookConnection) DeleteBook(book entity.Book) {
	db.connection.Delete(book)
}
func (db *bookConnection) AllBook() []entity.Book {
	var books []entity.Book
	db.connection.Preload("User").Find(&books)
	return books
}
func (db *bookConnection) FindBookByID(bookId uint64) entity.Book {
	var book entity.Book
	db.connection.Preload("User").Find(&book, bookId)
	return book
}
