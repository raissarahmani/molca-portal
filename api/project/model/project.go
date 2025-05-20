package model

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/molca-id/portal-app-api/arch/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongod "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const CollectionName = "projects"

type Project struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Title     string        `bson:"title" validate:"required,max=500"`
	Slug      string        `bson:"slug" validate:"required,max=200"`
	Link      string        `bson:"link" validate:"required,max=2000"`
	ImageURL  string        `bson:"image_url" validate:"required,max=2000"`
	Type      string        `bson:"type" validate:"required,oneof=digital-twin smart-manufacture ar vr deck tool"`
	CreatedAt time.Time     `bson:"created_at" validate:"required"`
	UpdatedAt time.Time     `bson:"updated_at" validate:"required"`
}

func NewProject(title, slug, link, imageURL, projectType string) (*Project, error) {
	now := time.Now()

	project := Project{
		Title:     title,
		Slug:      slug,
		Link:      link,
		ImageURL:  imageURL,
		Type:      projectType,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := project.Validate(); err != nil {
		return nil, err
	}

	return &project, nil
}

func (p *Project) GetValue() *Project {
	return p
}

func (p *Project) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (*Project) EnsureIndexes(db mongo.Database) {
	indexes := []mongod.IndexModel{
		{
			Keys: bson.D{
				{Key: "slug", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}

	mongo.NewQueryBuilder[Project](db, CollectionName).Query(context.Background()).CreateIndexes(indexes)
}
