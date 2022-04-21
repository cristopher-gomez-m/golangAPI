package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cristopher-gomez-m/golang_api/dto"
	"github.com/cristopher-gomez-m/golang_api/helper"
	"github.com/cristopher-gomez-m/golang_api/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserController interface {
	Update(context *gin.Context)
	Profile(context *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (controller *userController) Update(context *gin.Context) {
	var userUpdateDTO dto.UserUpdateDTO
	errDTO := context.ShouldBind(&userUpdateDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authHeader := context.GetHeader("Authorization")
	token, err := controller.jwtService.ValidateToken(authHeader)
	if err != nil {
		panic(err)
	}
	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err)
	}
	userUpdateDTO.ID = id
	update := controller.userService.Update(userUpdateDTO)
	response := helper.BuildResponse(true, "OK", update)
	context.JSON(http.StatusOK, response)

}
func (controller *userController) Profile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, err := controller.jwtService.ValidateToken(authHeader)
	if err != nil {
		panic(err.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := fmt.Sprintf("%v", claims["user_id"])
	user := controller.userService.Profile(userId)
	response := helper.BuildResponse(true, "OK", user)
	context.JSON(http.StatusOK, response)
}
