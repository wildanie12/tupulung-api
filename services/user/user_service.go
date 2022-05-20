package user

import (
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"
	"tupulung/deliveries/middleware"
	"tupulung/deliveries/validations"
	entity "tupulung/entities"
	"tupulung/entities/web"
	eventRepository "tupulung/repositories/event"
	userRepository "tupulung/repositories/user"
	storageProvider "tupulung/utilities/storage"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  userRepository.UserRepositoryInterface
	eventRepo eventRepository.EventRepositoryInterface
	validate  *validator.Validate
}

func NewUserService(repository userRepository.UserRepositoryInterface, eventRepo eventRepository.EventRepositoryInterface) *UserService {
	return &UserService{
		userRepo:  repository,
		eventRepo: eventRepo,
		validate:  validator.New(),
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
	userRes := entity.UserResponse{}
	copier.Copy(&userRes, &user)

	return userRes, err
}

/*
 * User Service - Get User joined events
 * -------------------------------
 * Mengambil data event yang sudah di join oleh user
 */
func (service UserService) GetJoinedEvents(userID int) ([]entity.EventResponse, error) {

	// Mengambil data user dari repository
	events, err := service.userRepo.GetJoinedEvents(userID)
	if err != nil {
		return []entity.EventResponse{}, err
	}

	// proses menjadi user response
	eventRes := []entity.EventResponse{}
	copier.Copy(&eventRes, &events)

	return eventRes, err
}

/*
 * User Service - Create (register)
 * -------------------------------
 * Registrasi User dan mengembalikan token dan auth response
 * untuk auto sign in setelah register
 */
func (service UserService) Create(userRequest entity.UserRequest, avatar *multipart.FileHeader, storageProvider storageProvider.StorageInterface) (entity.AuthResponse, error) {

	// Validation
	userFiles := []*multipart.FileHeader{}
	if avatar != nil {
		userFiles = append(userFiles, avatar)
	}
	err := validations.ValidateCreateUserRequest(service.validate, userRequest, userFiles)
	if err != nil {
		return entity.AuthResponse{}, err
	}

	// Konversi user request menjadi domain untuk diteruskan ke repository
	user := entity.User{}
	copier.Copy(&user, &userRequest)

	// Konversi datetime untuk field datetime (dob)
	dob, err := time.Parse("2006-01-02", userRequest.DOB)
	if err != nil {
		return entity.AuthResponse{}, web.WebError{Code: 400, Message: "date of birth format is invalid"}
	}
	user.DOB = dob

	// Password hashing menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.AuthResponse{}, web.WebError{Code: 500, Message: "server error: hashing failed"}
	}
	user.Password = string(hashedPassword)

	// Upload avatar if exists
	if avatar != nil {

		// Upload avatar to S3
		filename := uuid.New().String() + avatar.Filename
		avatarURL, err := storageProvider.UploadFromRequest("avatar/"+filename, avatar)
		if err != nil {
			return entity.AuthResponse{}, web.WebError{Code: 500, Message: err.Error()}
		}
		user.Avatar = avatarURL
	}

	// Insert ke sistem melewati repository
	user, err = service.userRepo.Store(user)
	if err != nil {
		return entity.AuthResponse{}, err
	}

	// Konversi hasil repository menjadi user response
	userRes := entity.UserResponse{}
	copier.Copy(&userRes, &user)

	// generate token
	token, err := middleware.CreateToken(user)
	if err != nil {
		return entity.AuthResponse{}, err
	}

	// Buat auth response untuk dimasukkan token dan user
	authRes := entity.AuthResponse{
		Token: token,
		User:  userRes,
	}
	return authRes, nil
}

/*
 * User Service - Update
 * -------------------------------
 * Edit data user / edit profile
 */
func (service UserService) Update(userRequest entity.UserRequest, userID int, avatar *multipart.FileHeader, storageProvider storageProvider.StorageInterface) (entity.UserResponse, error) {

	// validation
	userFiles := []*multipart.FileHeader{}
	if avatar != nil {
		userFiles = append(userFiles, avatar)
	}
	err := validations.ValidateUpdateUserRequest(userFiles)
	if err != nil {
		return entity.UserResponse{}, err
	}

	// Get user by ID via repository
	user, err := service.userRepo.Find(userID)
	if err != nil {
		return entity.UserResponse{}, web.WebError{Code: 400, Message: "The requested ID doesn't match with any record"}
	}

	// Hanya hash password jika password juga diganti (tidak kosong)
	if userRequest.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			return entity.UserResponse{}, web.WebError{Code: 500, Message: "server error: hashing failed"}
		}
		user.Password = string(hashedPassword)
	}

	// Avatar
	if avatar != nil {

		// Delete avatar lama jika ada yang baru
		if user.Avatar != "" {
			u, _ := url.Parse(user.Avatar)
			objectPathS3 := strings.TrimPrefix(u.Path, "/")
			storageProvider.Delete(objectPathS3)
		}

		// Upload avatar to S3
		filename := uuid.New().String() + avatar.Filename
		avatarURL, err := storageProvider.UploadFromRequest("avatar/"+filename, avatar)
		if err != nil {
			return entity.UserResponse{}, web.WebError{Code: 500, Message: err.Error()}
		}
		user.Avatar = avatarURL
	}

	// Konversi dari request ke domain entity user - mengabaikan nilai kosong pada request
	copier.CopyWithOption(&user, &userRequest, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	// Update via repository
	user, err = service.userRepo.Update(user, userID)

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
func (service UserService) Delete(userID int, storageProvider storageProvider.StorageInterface) error {

	// Cari user berdasarkan ID via repo
	user, err := service.userRepo.Find(userID)
	if err != nil {
		return web.WebError{Code: 400, Message: "The requested ID doesn't match with any record"}
	}

	// Delete avatar lama jika ada yang baru
	if user.Avatar != "" {
		u, _ := url.Parse(user.Avatar)
		objectPathS3 := strings.TrimPrefix(u.Path, "/")
		storageProvider.Delete(objectPathS3)
	}

	// Delete user event
	filters := []map[string]string{
		{
			"field":    "user_id",
			"operator": "=",
			"value":    strconv.Itoa(int(user.ID)),
		},
	}
	service.eventRepo.DeleteBatch(filters)

	// Delete via repository
	err = service.userRepo.Delete(userID)
	return err
}
