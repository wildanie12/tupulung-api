package validations

import (
	"mime/multipart"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
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
 * Filesize Validation Rules
 * -------------------------------
 * Aturan input file user berdasarkan size
 * [field]: [size]
 */
var userFileSizeRules = map[string]int{
	"avatar": 1024 * 1024,  		// 1MB 
}

/*
 * Filesize Validation Rules
 * -------------------------------
 * Aturan format ekstensi file yang diperbolehkan
 * [field]: ext1|ext2|ext3...
 */
var userFileExtRules = map[string]string{
	"avatar": "jpg|jpeg",
}

/*
 * User Validation - Validate Create User Request
 * -------------------------------
 * Validasi user saat registrasi berdasarkan validate tag 
 * yang ada pada user request dan file rules diatas
 */
func ValidateCreateUserRequest(validate *validator.Validate, userReq entities.UserRequest, userFiles []*multipart.FileHeader) error {

	errors := []web.ValidationErrorItem{}
	
	validateUserStruct(validate, userReq, &errors)
	validateUserFiles(userFiles, &errors)

	if len(errors) > 0 {
		return web.ValidationError{
			Code: 400,
			Message: "Validation error",
			Errors: errors,
		}
	}
	return nil
}

/*
 * User Validation - Validate Update User Request
 * -------------------------------
 * Validasi user saat edit profile berdasarkan
 * file rules diatas
 */
func ValidateUpdateUserRequest(userFiles []*multipart.FileHeader) error {

	errors := []web.ValidationErrorItem{}
	
	validateUserFiles(userFiles, &errors)
	if len(errors) > 0 {
		return web.ValidationError{
			Code: 400,
			Message: "Validation error",
			Errors: errors,
		}
	}
	return nil
}


func validateUserStruct(validate *validator.Validate, userReq entities.UserRequest, errors *[]web.ValidationErrorItem) {
	err := validate.Struct(userReq)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(userReq).FieldByName(err.Field())
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field.Tag.Get("form"),
				Error: userErrorMessages[err.Field() + "|" + err.ActualTag()],
			})
		}
	}
}



func validateUserFiles(userFiles []*multipart.FileHeader, errors *[]web.ValidationErrorItem) {
	// File validation
	for _, file := range userFiles {

		// Parse file header Content-Disposition to get field name
		field := strings.TrimPrefix(strings.Split(file.Header.Get("Content-Disposition"), ";")[1], " ")
		field = strings.Split(field, "=")[1]
		field = strings.Replace(field, "\"", "", -1)
		field = strings.Replace(field, "\\", "", -1)

		
		// Size validations
		if file.Size > int64(userFileSizeRules[field]) {
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field,
				Error: field + " size cannot more than " + strconv.Itoa(userFileSizeRules[field] / 1024) + " KB",
			})
		}

		// Extension validations
		fileExt := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		allowedExt := strings.Split(userFileExtRules[field], "|")
		fileExtAllowed := false
		for _, ext := range allowedExt {
			if fileExt == ext {
				fileExtAllowed = true
				break
			}
		}
		if !fileExtAllowed {
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field,
				Error: field + " field must be type of " + strings.Join(allowedExt, ", "),
			})
		}
	}
}