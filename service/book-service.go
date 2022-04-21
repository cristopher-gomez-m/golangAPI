package service

import (
	"fmt"
	"log"

	"github.com/cristopher-gomez-m/golang_api/dto"
	"github.com/cristopher-gomez-m/golang_api/entity"
	"github.com/cristopher-gomez-m/golang_api/repository"
	"github.com/mashingan/smapping"
)

type BookService interface {
	Insert(book dto.BookCreateDTO) entity.Book
	Update(book dto.BookUpdateDTO) entity.Book
	Delete(book entity.Book)
	All() []entity.Book
	FindById(bookId uint64) entity.Book
	IsAllowedToEdit(userId string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRepo,
	}
}

func (service *bookService) Insert(book dto.BookCreateDTO) entity.Book {
	bookToInsert := entity.Book{}
	err := smapping.FillStruct(&bookToInsert, smapping.MapFields(&book))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	response := service.bookRepository.InsertBook(bookToInsert)
	return response
}
func (service *bookService) Update(book dto.BookUpdateDTO) entity.Book {
	bookToInsert := entity.Book{}
	err := smapping.FillStruct(&bookToInsert, smapping.MapFields(&book))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	response := service.bookRepository.UpdateBook(bookToInsert)
	return response
}
func (service *bookService) Delete(book entity.Book) {
	service.bookRepository.DeleteBook(book)
}
func (service *bookService) All() []entity.Book {
	return service.bookRepository.AllBook()
}
func (service *bookService) FindById(bookId uint64) entity.Book {
	return service.bookRepository.FindBookByID(bookId)
}
func (service *bookService) IsAllowedToEdit(userId string, bookID uint64) bool {
	book := service.bookRepository.FindBookByID(bookID)
	id := fmt.Sprintf("%v", book.User.ID)
	return userId == id
}
