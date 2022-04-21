package controllers

import (
	"net/http"
	"strconv"

	"github.com/cristopher-gomez-m/golang_api/dto"
	"github.com/cristopher-gomez-m/golang_api/entity"
	"github.com/cristopher-gomez-m/golang_api/helper"
	"github.com/cristopher-gomez-m/golang_api/service"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(context *gin.Context)
	Register(context *gin.Context)
}

type authController struct {
	authService service.AuthService
	jwtService  service.JWTService
}

func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(context *gin.Context) {
	var loginDto dto.LoginDTO
	errDTO := context.ShouldBind(&loginDto)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authResult := c.authService.VerifyCredential(loginDto.Email, loginDto.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10))
		v.Token = generatedToken
		response := helper.BuildResponse(true, "OK", v)
		context.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorResponse("Please check again your credential", "Invalid credential", helper.EmptyObj{})
	context.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (c *authController) Register(context *gin.Context) {
	var registerDTO dto.RegisterDTO
	errDTO := context.ShouldBind(&registerDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	if c.authService.IsDuplicatedEmail(registerDTO.Email) {
		response := helper.BuildErrorResponse("Failed to process request", "Duplicate email", helper.EmptyObj{})
		context.JSON(http.StatusConflict, response)
	} else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		response := helper.BuildResponse(true, "OK", createdUser)
		context.JSON(http.StatusCreated, response)
	}
}
