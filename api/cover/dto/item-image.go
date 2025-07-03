package dto

import (
	"github.com/go-playground/validator/v10"
)

type ItemImage struct {
	ImageUrl string `json:"image_url" binding:"required"`
}

func EmptyItemImage() *ItemImage {
	return &ItemImage{}
}

func (i *ItemImage) GetValue() *ItemImage {
	return i
}

func (d *ItemImage) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, err.Field()+" is required")
		case "min":
			msgs = append(msgs, err.Field()+" must be at least 3 characters long")
		case "max":
			msgs = append(msgs, err.Field()+" must be at most 255 characters long")
		default:
			msgs = append(msgs, err.Field()+" is invalid")
		}
	}

	return msgs, nil
}
