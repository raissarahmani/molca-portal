package dto

import (
	"github.com/go-playground/validator/v10"
)

type FilterPaginatedProject struct {
	Search *string `form:"search,omitempty" validate:"omitempty"`
	Type   *string `form:"type,omitempty" validate:"omitempty,oneof=digital-twin smart-manufacture ar vr deck tool"`
	Page   int64   `form:"page" binding:"required" validate:"min=1,max=1000"`
	Limit  int64   `form:"limit" binding:"required" validate:"min=1,max=1000"`
}

func (f *FilterPaginatedProject) GetValue() *FilterPaginatedProject {
	return f
}

func EmptyFilterPaginatedProject() *FilterPaginatedProject {
	return &FilterPaginatedProject{}
}

func (d *FilterPaginatedProject) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, err.Field()+" is required")
		case "min":
			msgs = append(msgs, err.Field()+" must be at least 3 characters long")
		case "max":
			msgs = append(msgs, err.Field()+" must be at most 255 characters long")
		case "oneof":
			msgs = append(msgs, err.Field()+" must be one of the following: digital-twin, smart-manufacture, ar, vr, deck, tool")
		}
	}

	return msgs, nil
}
