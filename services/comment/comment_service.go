package comment

import (
	"reflect"
	"tupulung/entities"
	"tupulung/entities/web"
	commentRepo "tupulung/repositories/comment"
	userRepo "tupulung/repositories/user"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
)

type CommentService struct {
	commentRepo commentRepo.CommentRepositoryInterface
	userRepo userRepo.UserRepositoryInterface
}

func NewCommentService(commentRepo commentRepo.CommentRepositoryInterface, userRepo userRepo.UserRepositoryInterface) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		userRepo: userRepo,
	}
}

/*
 * Find All
 * -------------------------------
 * Mengambil data comment berdasarkan filters dan sorts
 */
func (service CommentService) FindAll(limit, page int, filters []map[string]string, sorts []map[string]interface{}) ([]entities.CommentResponse, error) {

	offset := (page - 1) * limit

	// Repository action find all comment
	comments, err := service.commentRepo.FindAll(limit, offset, filters, sorts)
	if err != nil {
		return []entities.CommentResponse{}, err
	}

	// Konversi ke comment response
	commentsRes := []entities.CommentResponse{}
	copier.Copy(&commentsRes, &comments)
	return commentsRes, nil
}

/*
 * Load pagination data
 * -------------------------------
 * Mengambil data pagination comment berdasarkan filters
 */
func (service CommentService) GetPagination(page, limit int, filters []map[string]string) (web.Pagination, error) {
	totalRows, err := service.commentRepo.CountAll(filters)
	if err != nil {
		return web.Pagination{}, err
	}
	var totalPages int64 = 1
	if limit > 0 {
		totalPages = totalRows / int64(limit)
	}
	if totalPages <= 0 {
		totalPages = 1
	}
	return web.Pagination{
		Page: page,
		Limit: limit,
		TotalPages: int(totalPages),
	}, nil
}

/*
 * Create comments
 * -------------------------------
 * Membuat komentar baru berdasarkan user yang sedang login
 */
func (service CommentService) Create(commentRequest entities.CommentRequest, eventID int, tokenReq interface{}) (entities.CommentResponse, error) {
	
	// convert request to domain entity
	comment := entities.Comment{}
	copier.Copy(&comment, &commentRequest)

	// Translate token
	token := tokenReq.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDReflect := reflect.ValueOf(claims).MapIndex(reflect.ValueOf("userID"))
	if reflect.ValueOf(userIDReflect.Interface()).Kind().String() != "float64" {
		return entities.CommentResponse{}, web.WebError{Code: 400, Message: "Invalid token, no userdata present"}
	}

	// get user data
	user, err := service.userRepo.Find(int(claims["userID"].(float64)))
	if err != nil {
		return entities.CommentResponse{}, web.WebError{Code: 400, Message: "No user matched with this authenticated user"}
	}
	comment.UserID = user.ID
	comment.EventID = uint(eventID)

	// Repository action
	comment, err = service.commentRepo.Store(comment)
	if err != nil {
		return entities.CommentResponse{}, err
	}

	// process domain entity to response
	commentRes := entities.CommentResponse{}
	copier.Copy(&commentRes, &comment)

	return commentRes, nil
}

/*
 * Update Comment
 * -------------------------------
 * Edit komentar user, hanya pemilik komentar yang dapat mengedit
 */
func (service CommentService) Update(commentRequest entities.CommentRequest, id int, tokenReq interface{}) (entities.CommentResponse, error) {

	// Find comment
	comment, err := service.commentRepo.Find(id)
	if err != nil {
		return entities.CommentResponse{}, web.WebError{ Code: 400, Message: "The requested ID doesn't match with any record" }
	}

	// Translate token
	token := tokenReq.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDReflect := reflect.ValueOf(claims).MapIndex(reflect.ValueOf("userID"))
	if reflect.ValueOf(userIDReflect.Interface()).Kind().String() != "float64" {
		return entities.CommentResponse{}, web.WebError{Code: 400, Message: "Invalid token, no userdata present"}
	}

	// Match comment with authenticated userid
	userID := int(claims["userID"].(float64))
	if userID != int(comment.UserID) {
		return entities.CommentResponse{}, web.WebError{ Code: 401, Message: "Unauthorized user, cannot update someone else's comment" }
	}

	// Merge updated data request to domain entity
	copier.CopyWithOption(&comment, &commentRequest, copier.Option{ IgnoreEmpty: true, DeepCopy: true })

	// repository action
	comment, err = service.commentRepo.Update(comment, id)

	// Convert comment domain to comment response
	commentRes := entities.CommentResponse{}
	copier.Copy(&commentRes, &comment)

	return commentRes, err
}

/*
 * Delete Comment
 * -------------------------------
 * Hapus komentar user, hanya pemilik komentar yang dapat mengedit
 */
func (service CommentService) Delete(id int, tokenReq interface{}) error {
	// Find comment
	comment, err := service.commentRepo.Find(id)
	if err != nil {
		return web.WebError{ Code: 400, Message: "The requested ID doesn't match with any record" }
	}

	// Translate token
	token := tokenReq.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDReflect := reflect.ValueOf(claims).MapIndex(reflect.ValueOf("userID"))
	if reflect.ValueOf(userIDReflect.Interface()).Kind().String() != "float64" {
		return web.WebError{Code: 400, Message: "Invalid token, no userdata present"}
	}

	// Match comment with authenticated userid
	userID := int(claims["userID"].(float64))
	if userID != int(comment.UserID) {
		return web.WebError{ Code: 401, Message: "Unauthorized user, cannot Delete someone else's comment" }
	}
	
	// Copy request to found comment
	err = service.commentRepo.Delete(id)
	return err
}