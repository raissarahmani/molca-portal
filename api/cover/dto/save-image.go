package dto

import (
	"fmt"
	"mime/multipart"

	"github.com/go-playground/validator/v10"
)

type SaveImage struct {
	Name string                `form:"name" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required" validate:"required"`
}

func EmptySaveImage() *SaveImage {
	return &SaveImage{}
}

func (si *SaveImage) GetValue() *SaveImage {
	return si
}

func (si *SaveImage) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string

	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s is too long", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s is too short", err.Field()))
		case "mimetypes":
			msgs = append(msgs, fmt.Sprintf("%s is not a valid image", err.Field()))
		case "maxfilesize":
			msgs = append(msgs, fmt.Sprintf("%s is too large", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}

	return msgs, nil
}
