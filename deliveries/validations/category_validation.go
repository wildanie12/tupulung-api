package validations

import (
	"fmt"
	"reflect"
	"tupulung/entities"
	"tupulung/entities/web"
	"tupulung/utilities"

	"github.com/go-playground/validator/v10"
)

/*
 * Category Validation - Error Message
 * -------------------------------
 * Kumpulan custom error message yang ditampilkan
 * ke response berdasarkan struct field dan validate tagnya
 */
var categoryErrorMessages = map[string]string {
	"Title|required": "Title field must be filled",
}

/*
 * Category Validation - Validate
 * -------------------------------
 * Validasi category berdasarkan validate tag 
 * yang ada pada category request
 */
func ValidateCategoryRequest(validate *validator.Validate, categoryReq entities.CategoryRequest) error {
	err := validate.Struct(categoryReq)
	fmt.Println(utilities.JsonEncode(categoryReq))
	if err != nil {
		errors := []web.ValidationErrorItem{}
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(categoryReq).FieldByName(err.Field())
			errors = append(errors, web.ValidationErrorItem{
				Field: field.Tag.Get("form"),
				Error: categoryErrorMessages[err.Field() + "|" + err.ActualTag()],
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