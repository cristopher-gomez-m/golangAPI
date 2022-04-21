package service

import (
	"log"

	"github.com/cristopher-gomez-m/golang_api/dto"
	"github.com/cristopher-gomez-m/golang_api/entity"
	"github.com/cristopher-gomez-m/golang_api/repository"
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredential(email string, password string) interface{}
	CreateUser(user dto.RegisterDTO) entity.User
	FindByEmail(email string) entity.User
	IsDuplicatedEmail(email string) bool
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepo,
	}
}

func (service *authService) VerifyCredential(email string, password string) interface{} {
	res := service.userRepository.VerifyCredential(email, password)
	if v, ok := res.(entity.User); ok {
		comparedPassword := comparePassword(v.Password, []byte(password))
		if v.Email == email && comparedPassword {
			return res
		}
		return false
	}
	return false
}

func comparePassword(hashedPWD string, plainPaswword []byte) bool {
	byteHash := []byte(hashedPWD)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPaswword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
func (service *authService) CreateUser(user dto.RegisterDTO) entity.User {
	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	res := service.userRepository.InsertUser(userToCreate)
	return res
}
func (service *authService) FindByEmail(email string) entity.User {
	return service.userRepository.FindByEmail(email)
}
func (service *authService) IsDuplicatedEmail(email string) bool {
	res := service.userRepository.IsDuplicatedEmail(email)
	return (res.Error == nil)
}
