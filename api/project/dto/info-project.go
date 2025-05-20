package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type InfoProject struct {
	ID        bson.ObjectID `json:"_id" binding:"required"`
	Title     string        `json:"title" binding:"required"`
	Slug      string        `json:"slug" binding:"required"`
	Link      string        `json:"link" binding:"required"`
	ImageUrl  string        `json:"image_url" binding:"required"`
	Type      string        `json:"type" binding:"required"`
	CreatedAt time.Time     `json:"created_at" binding:"required"`
	UpdatedAt time.Time     `json:"last_updated_at" binding:"required"`
}

func EmptyInfoProject() *InfoProject {
	return &InfoProject{}
}

func (d *InfoProject) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
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
