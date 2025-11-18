package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidatorError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func BindValidate[T any](c *gin.Context, req *T) ([]ValidatorError, error) {
	if err := c.ShouldBindJSON(req); err != nil {
		var errs []ValidatorError
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				errs = append(errs, ValidatorError{
					Field:   e.Field(),
					Message: validationMessage(e),
				})
			}
			return errs, err
		}
	}
	return nil, nil
}

func validationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return "must be at least " + e.Param() + " characters"
	case "max":
		return "must be at most " + e.Param() + " characters"
	case "gte":
		return "must be greater than or equal to " + e.Param()
	case "lte":
		return "must be less than or equal to " + e.Param()
	default:
		return "is not valid"
	}
}
