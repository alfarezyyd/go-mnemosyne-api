package config

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

func InitializeValidator() (*validator.Validate, ut.Translator) {
	var validatorInstance = validator.New()
	// Tambahkan validasi kustom untuk 'datetime'
	validatorInstance.RegisterValidation("date", func(fieldLevel validator.FieldLevel) bool {
		format := "2006-01-02"
		_, err := time.Parse(format, fieldLevel.Field().String())
		return err == nil
	})
	validatorInstance.RegisterValidation("maxSize", maxFileSizeValidation)
	validatorInstance.RegisterValidation("extensionFile", validateFileExtensionValidation)
	validatorInstance.RegisterValidation("obligatoryFile", requiredFileValidationValidation)

	englishLang := en.New()
	universalTranslator := ut.New(englishLang, englishLang)
	engTranslator, _ := universalTranslator.GetTranslator("en")
	_ = translations.RegisterDefaultTranslations(validatorInstance, engTranslator)

	return validatorInstance, engTranslator
}

func requiredFileValidationValidation(fieldLevel validator.FieldLevel) bool {
	_, ok := fieldLevel.Field().Interface().(multipart.FileHeader)
	if !ok {
		return false
	}
	return true
}

func maxFileSizeValidation(fieldLevel validator.FieldLevel) bool {
	file, ok := fieldLevel.Field().Interface().(multipart.FileHeader)
	if !ok {
		return false
	}

	// Dapatkan parameter maxSize
	maxSize, err := strconv.ParseInt(fieldLevel.Param(), 10, 64)
	if err != nil {
		return false
	}

	// Validasi ukuran file
	return file.Size <= maxSize*1024*1024
}

// Validasi ekstensi file
func validateFileExtensionValidation(fieldLevel validator.FieldLevel) bool {
	// Ambil file dari field
	file, ok := fieldLevel.Field().Interface().(multipart.FileHeader)
	if !ok {
		return false
	}

	// Ambil parameter ekstensi yang diperbolehkan
	allowedExtension := fieldLevel.Param() // Misalnya "png,jpg,jpeg"

	// Pisahkan parameter ekstensi menjadi slice
	extSlice := strings.Split(allowedExtension, " ")

	// Ambil ekstensi file
	filename := strings.ToLower(file.Filename)
	for _, ext := range extSlice {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}
