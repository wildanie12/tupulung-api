package user

import (
	"reflect"
	"time"
	entity "tupulung/entities"
	"tupulung/entities/web"
	userRepository "tupulung/repositories/user"
	"tupulung/utilities"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo userRepository.UserRepositoryInterface
}

func NewUserService(repository userRepository.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepo: repository,
	}
}

/*
 * User Service - Find
 * -------------------------------
 * Mencari user berdasarkan ID
 */
func (service UserService) Find(id int) (entity.UserResponse, error) {
	
	// Mengambil data user dari repository
	user, err := service.userRepo.Find(id)

	// proses menjadi user response
	userRes  := entity.UserResponse{}
	copier.Copy(&userRes, &user)

	return userRes, err
}

/*
 * User Service - Create (register)
 * -------------------------------
 * Registrasi User dan mengembalikan token dan auth response
 * untuk auto sign in setelah register
 */
func (service UserService) Create(userRequest entity.UserRequest) (entity.AuthResponse, error) {
	
	// Konversi user request menjadi domain untuk diteruskan ke repository 
	user := entity.User{}
	copier.Copy(&user, &userRequest)

	// Konversi datetime untuk field datetime (dob)
	dob, err := time.Parse("2006-01-02", userRequest.DOB)
	if err != nil {
		return entity.AuthResponse{}, web.WebError{ Code: 400, Message: "date of birth format is invalid" }
	}
	user.DOB = dob


	// Password hashing menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.AuthResponse{}, web.WebError{ Code: 500, Message: "server error: hashing failed" }
	}
	user.Password = string(hashedPassword)

	// Insert ke sistem melewati repository
	user, err = service.userRepo.Store(user)
	if err != nil {
		return entity.AuthResponse{}, err
	}

	// Konversi hasil repository menjadi user response
	userRes := entity.UserResponse{}
	copier.Copy(&userRes, &user)

	// generate token
	token, err := utilities.CreateToken(user)
	if err != nil {
		return entity.AuthResponse{}, err
	}

	// Buat auth response untuk dimasukkan token dan user
	authRes := entity.AuthResponse{
		Token: token,
		User: userRes,
	}
	return authRes, nil
}


/*
 * User Service - Update 
 * -------------------------------
 * Edit data user / edit profile
 */
func (service UserService) Update(userRequest entity.UserRequest, id int, tokenReq interface{}) (entity.UserResponse, error) {

	// Translate token
	token := tokenReq.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDReflect := reflect.ValueOf(claims).MapIndex(reflect.ValueOf("userID"))  // Mengambil tipe data dari claims userID
	if reflect.ValueOf(userIDReflect.Interface()).Kind().String() != "float64" { // Tolak jika bukan float64
		return entity.UserResponse{}, web.WebError{ Code: 400, Message: "Invalid token, no userdata present" }
	}

	// Reject jika id user yg dirubah tidak sama dengan id user token
	if id != int(claims["userID"].(float64)) {
		return entity.UserResponse{}, web.WebError{ Code: 401, Message: "Unauthorized user" }
	}

	// Get user by ID via repository
	user, err := service.userRepo.Find(id)
	if err != nil {
		return entity.UserResponse{}, web.WebError{ Code: 400, Message: "The requested ID doesn't match with any record" }
	}
	
	// Konversi dari request ke domain entity user - mengabaikan nilai kosong pada request
	copier.CopyWithOption(&user, &userRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	// Hanya hash password jika password juga diganti (tidak kosong)
	if userRequest.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			return entity.UserResponse{}, web.WebError{ Code: 500, Message: "server error: hashing failed" }
		}
		user.Password = string(hashedPassword)
	}

	// Update via repository
	user, err = service.userRepo.Update(user, id)

	// Konversi user domain menjadi user response
	userRes := entity.UserResponse{}
	copier.Copy(&userRes, &user)

	return userRes, err
}

/*
 * User Service - Delete
 * -------------------------------
 * Delete data user menggunakan ID
 * Hanya usernya sendiri yang dapat melakukan delete
 */
func (service UserService) Delete(id int, tokenReq interface{}) error {

	// Translate token
	token := tokenReq.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDReflect := reflect.ValueOf(claims).MapIndex(reflect.ValueOf("userID"))  // Mengambil tipe data dari claims userID
	if reflect.ValueOf(userIDReflect.Interface()).Kind().String() != "float64" { // Tolak jika bukan float64
		return web.WebError{ Code: 400, Message: "Invalid token, no userdata present" }
	}

	// Reject jika id user yg dirubah tidak sama dengan id user token
	if id != int(claims["userID"].(float64)) {
		return web.WebError{ Code: 401, Message: "Unauthorized user" }
	}

	// Cari user berdasarkan ID via repo
	_, err := service.userRepo.Find(id)
	if err != nil {
		return web.WebError{ Code: 400, Message: "The requested ID doesn't match with any record" }
	}
	
	// Delete via repository
	err = service.userRepo.Delete(id)
	return err
}