package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CreateProject struct {
	Title    string `json:"title" validate:"required,max=500"`
	Slug     string `json:"slug" validate:"required,max=200"`
	Link     string `json:"link" validate:"required,max=2000"`
	ImageURL string `json:"image_url" validate:"required,max=2000"`
	Type     string `json:"type" validate:"required,oneof=digital-twin smart-manufacture ar vr deck tool"`
}

func EmptyCreateProject() *CreateProject {
	return &CreateProject{}
}

func (c *CreateProject) GetValue() *CreateProject {
	return c
}

func (c *CreateProject) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string

	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s is too long", err.Field()))
		case "oneof":
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}

	return msgs, nil
}
