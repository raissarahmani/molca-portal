package coredto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func EmptyPagination() *Pagination {
	return &Pagination{}
}

type Pagination struct {
	Page  int64 `form:"page" binding:"required" validate:"min=1,max=1000"`
	Limit int64 `form:"limit" binding:"required" validate:"min=1,max=1000"`
}

func (p *Pagination) GetValue() *Pagination {
	return p
}

func (p *Pagination) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s must be less than %s", err.Field(), err.Param()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
