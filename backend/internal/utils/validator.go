package utils

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validators
	validate.RegisterValidation("alphanum_hyphen", validateAlphanumHyphen)
	validate.RegisterValidation("password_strength", validatePasswordStrength)
}

// ValidateStruct validates a struct using validator tags
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// validateAlphanumHyphen validates that a string contains only lowercase alphanumeric characters and hyphens
func validateAlphanumHyphen(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, value)
	return matched
}

// validatePasswordStrength validates password strength
func validatePasswordStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// GetValidationErrors returns a map of field errors from validation error
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			tag := fieldError.Tag()

			switch tag {
			case "required":
				errors[field] = field + " is required"
			case "email":
				errors[field] = "Invalid email format"
			case "min":
				errors[field] = field + " must be at least " + fieldError.Param() + " characters"
			case "max":
				errors[field] = field + " must be at most " + fieldError.Param() + " characters"
			case "alphanum_hyphen":
				errors[field] = field + " must contain only lowercase letters, numbers, and hyphens"
			case "password_strength":
				errors[field] = "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character"
			default:
				errors[field] = field + " validation failed"
			}
		}
	}

	return errors
}
