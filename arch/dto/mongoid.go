package coredto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/molca-id/portal-app-api/arch/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func EmptyMongoId() *MongoId {
	return &MongoId{}
}

type MongoId struct {
	Id string        `uri:"id" binding:"required" validate:"required,len=24"`
	ID bson.ObjectID `uri:"-" validate:"-"`
}

func (d *MongoId) GetValue() *MongoId {
	id, err := mongo.NewObjectID(d.Id)
	if err == nil {
		d.ID = id
	}

	return d
}

func (m *MongoId) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "len":
			msgs = append(msgs, fmt.Sprintf("%s must be 24 characters long", err.Field()))
		}
	}
	return msgs, nil
}
