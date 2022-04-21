package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cristopher-gomez-m/golang_api/dto"
	"github.com/cristopher-gomez-m/golang_api/entity"
	"github.com/cristopher-gomez-m/golang_api/helper"
	"github.com/cristopher-gomez-m/golang_api/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type BookController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type bookController struct {
	bookService service.BookService
	jwtService  service.JWTService
}

func NewBookController(bookService service.BookService, jwtService service.JWTService) BookController {
	return &bookController{
		bookService: bookService,
		jwtService:  jwtService,
	}
}

func (controller *bookController) getUserIDByToken(token string) string {
	validatedToken, err := controller.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := validatedToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}

func (controller *bookController) All(context *gin.Context) {
	var books []entity.Book = controller.bookService.All()
	response := helper.BuildResponse(true, "Ok", books)
	context.JSON(http.StatusOK, response)
}
func (controller *bookController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	var book entity.Book = controller.bookService.FindById(id)
	if (book == entity.Book{}) {
		response := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, response)
	} else {
		response := helper.BuildResponse(true, "Ok", book)
		context.JSON(http.StatusOK, response)
	}

}
func (controller *bookController) Insert(context *gin.Context) {
	var bookCreateDTO dto.BookCreateDTO
	errDTO := context.ShouldBind(&bookCreateDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	} else {
		authHeader := context.GetHeader("Authorization")
		userID := controller.getUserIDByToken(authHeader)
		convertedUserId, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			bookCreateDTO.UserId = convertedUserId
		}
		result := controller.bookService.Insert(bookCreateDTO)
		response := helper.BuildResponse(true, "Ok", result)
		context.JSON(http.StatusCreated, response)
	}
}
func (controller *bookController) Update(context *gin.Context) {
	var bookUpdateDTO dto.BookUpdateDTO
	errDTO := context.ShouldBind(&bookUpdateDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	authHeader := context.GetHeader("Authorization")
	token, errToken := controller.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if controller.bookService.IsAllowedToEdit(userID, bookUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			bookUpdateDTO.UserId = id
		}
		result := controller.bookService.Update(bookUpdateDTO)
		response := helper.BuildResponse(true, "Ok", result)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}
func (controller *bookController) Delete(context *gin.Context) {
	var book entity.Book
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to get id", "No param id were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	book.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := controller.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if controller.bookService.IsAllowedToEdit(userID, book.ID) {
		controller.bookService.Delete(book)
		response := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}
