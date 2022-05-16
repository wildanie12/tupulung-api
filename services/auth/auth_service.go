package auth

import (
	"tupulung/deliveries/middleware"
	userRepository "tupulung/repositories/user"

	"tupulung/entities"
	web "tupulung/entities/web"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo userRepository.UserRepositoryInterface
}

func NewAuthService(userRepo userRepository.UserRepositoryInterface) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

/*
 * Auth Service - Login
 * -------------------------------
 * Mencari user berdasarkan ID
 */
func (service AuthService) Login(authReq entities.AuthRequest) (entities.AuthResponse, error) {

	// Get user by username via repository
	user, err := service.userRepo.FindBy("email", authReq.Email)
	if err != nil {
		return entities.AuthResponse{}, web.WebError{Code: 401, Message: "Invalid credential"}
	}

	// Verify password
	match := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authReq.Password))
	if match != nil {
		return entities.AuthResponse{}, web.WebError{Code: 401, Message: "Invalid password"}
	}

	// Konversi menjadi user response
	userRes := entities.UserResponse{}
	copier.Copy(&userRes, &user)

	// Create token
	token, err := middleware.CreateToken(user)
	if err != nil {
		return entities.AuthResponse{}, web.WebError{Code: 500, Message: "Error create token"}
	}

	return entities.AuthResponse{
		Token: token,
		User:  userRes,
	}, nil
}

/*
 * Auth Service - Me
 * -------------------------------
 * Mendapatkan userdata yang sedang login
 */
func (service AuthService) Me(ID int, token interface{}) (interface{}, error) {

	userJWT := token.(*jwt.Token)
	// Get userdata via repository
	user, err := service.userRepo.Find(ID)
	userRes := entities.UserResponse{}
	copier.Copy(&userRes, &user)

	// Bentuk auth response
	authRes := entities.AuthResponse{
		Token: userJWT.Raw,
		User:  userRes,
	}

	return authRes, err
}
