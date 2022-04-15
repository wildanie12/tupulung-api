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
 * Event Validation - Error Message
 * -------------------------------
 * Kumpulan custom error message yang ditampilkan
 * ke response berdasarkan struct field dan validate tagnya
 */
var eventErrorMessages = map[string]string{
	"Title|required":         "Title field must be filled",
	"HostedBy|required":      "HostedBy field must be filled",
	"CategryID|required":     "CategryID field must be filled",
	"DatetimeEvent|required": "DatetimeEvent field must be filled",
	"Location|required":      "Location field must be filled",
	"Description|required":   "Description field must be filled",
}

/*
 * Filesize Validation Rules
 * -------------------------------
 * Aturan input file event berdasarkan size
 * [field]: [size]
 */
var eventFileSizeRules = map[string]int{
	"cover": 2024 * 2024, // 2MB
}

/*
 * Filesize Validation Rules
 * -------------------------------
 * Aturan format ekstensi file yang diperbolehkan
 * [field]: ext1|ext2|ext3...
 */
var eventFileExtRules = map[string]string{
	"cover": "jpg|jpeg|png|webp|bmp",
}

/*
 * event Validation - Validate Create event Request
 * -------------------------------
 * Validasi event saat registrasi berdasarkan validate tag
 * yang ada pada event request dan file rules diatas
 */
func ValidateCreateEventRequest(validate *validator.Validate, eventReq entities.EventRequest, eventFiles []*multipart.FileHeader) error {

	errors := []web.ValidationErrorItem{}

	validateEventStruct(validate, eventReq, &errors)
	validateEventFiles(eventFiles, &errors)

	if len(errors) > 0 {
		return web.ValidationError{
			Code:    400,
			Message: "Validation error",
			Errors:  errors,
		}
	}
	return nil
}

/*
 * Event Validation - Validate Update Event Request
 * -------------------------------
 * Validasi event saat edit event berdasarkan
 * file rules diatas
 */
func ValidateUpdateEventRequest(eventFiles []*multipart.FileHeader) error {

	errors := []web.ValidationErrorItem{}

	validateEventFiles(eventFiles, &errors)
	if len(errors) > 0 {
		return web.ValidationError{
			Code:    400,
			Message: "Validation error",
			Errors:  errors,
		}
	}
	return nil
}

func validateEventStruct(validate *validator.Validate, eventReq entities.EventRequest, errors *[]web.ValidationErrorItem) {
	err := validate.Struct(eventReq)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(eventReq).FieldByName(err.Field())
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field.Tag.Get("form"),
				Error: eventErrorMessages[err.Field()+"|"+err.ActualTag()],
			})
		}
	}
}

func validateEventFiles(eventFiles []*multipart.FileHeader, errors *[]web.ValidationErrorItem) {
	// File validation
	for _, file := range eventFiles {

		// Parse file header Content-Disposition to get field name
		field := strings.TrimPrefix(strings.Split(file.Header.Get("Content-Disposition"), ";")[1], " ")
		field = strings.Split(field, "=")[1]
		field = strings.Replace(field, "\"", "", -1)
		field = strings.Replace(field, "\\", "", -1)

		// Size validations
		if file.Size > int64(eventFileSizeRules[field]) {
			*errors = append(*errors, web.ValidationErrorItem{
				Field: field,
				Error: field + " size cannot more than " + strconv.Itoa(eventFileSizeRules[field]/2024) + " KB",
			})
		}

		// Extension validations
		fileExt := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		allowedExt := strings.Split(eventFileExtRules[field], "|")
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
