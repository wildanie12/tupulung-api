package validations

import (
	"reflect"
	"tupulung/entities"
	"tupulung/entities/web"

	"github.com/go-playground/validator/v10"
)

/*
 * User Validation - Error Message
 * -------------------------------
 * Kumpulan custom error message yang ditampilkan
 * ke response berdasarkan struct field dan validate tagnya
 */
var userErrorMessages = map[string]string {
	"Name|required": "Name field must be filled",
	"Email|required": "Email field must be filled",
	"Email|email": "Email field is not an email",
	"Password|required": "Password field must be filled",
	"Gender|required": "Gender field must be filled",
	"DOB|required": "Date of birth field must be filled",
}

/*
 * User Validation - Validate
 * -------------------------------
 * Validasi user berdasarkan validate tag 
 * yang ada pada user request
 */
func ValidateUserRequest(validate *validator.Validate, userReq entities.UserRequest) error {
	err := validate.Struct(userReq)
	if err != nil {
		errors := []web.ValidationErrorItem{}
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(userReq).FieldByName(err.Field())
			errors = append(errors, web.ValidationErrorItem{
				Field: field.Tag.Get("form"),
				Error: userErrorMessages[err.Field() + "|" + err.ActualTag()],
			})
		}
		return web.ValidationError{
			Code: 400,
			Message: "Validation error",
			Errors: errors,
		}
	}
	return nil
}