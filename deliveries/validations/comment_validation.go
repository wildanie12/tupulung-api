package validations

import (
	"reflect"
	"tupulung/entities"
	"tupulung/entities/web"

	"github.com/go-playground/validator/v10"
)

/*
 * Comment Validation - Error Message
 * -------------------------------
 * Kumpulan custom error message yang ditampilkan
 * ke response berdasarkan struct field dan validate tagnya
 */
var commentErrorMessages = map[string]string {
	"Comment|required": "Comment field must be filled",
}
/*
 * Comment Validation - Validate Create Comment Request
 * -------------------------------
 * Validasi comment saat registrasi berdasarkan validate tag 
 * yang ada pada comment request dan file rules diatas
 */
func ValidateCreateCommentRequest(validate *validator.Validate, commentReq entities.CommentRequest) error {

	errors := []web.ValidationErrorItem{}
	
	validateCommentStruct(validate, commentReq, &errors)

	if len(errors) > 0 {
		return web.ValidationError{
			Code: 400,
			Message: "Validation error",
			Errors: errors,
		}
	}
	return nil
}



func validateCommentStruct(validate *validator.Validate, commentReq entities.CommentRequest, errors *[]web.ValidationErrorItem) {
	err := validate.Struct(commentReq)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(commentReq).FieldByName(err.Field())
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field.Tag.Get("form"),
				Error: commentErrorMessages[err.Field() + "|" + err.ActualTag()],
			})
		}
	}
}
