package project

import (
	"time"

	"github.com/molca-id/portal-app-api/api/project/dto"
	"github.com/molca-id/portal-app-api/api/project/model"

	"github.com/molca-id/portal-app-api/arch/mongo"
	"github.com/molca-id/portal-app-api/arch/network"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Service interface {
	SaveProject(d *dto.CreateProject) (*model.Project, error)
	GetProjectById(id bson.ObjectID) (*model.Project, error)
	GetProjectBySlug(slug string) (*model.Project, error)
	UpdateProject(id bson.ObjectID, d *dto.UpdateProject) error
	GetProjectSlugExists(slug string) bool
	DeleteProject(id bson.ObjectID) error
}

type service struct {
	network.BaseService
	projectQueryBuilder mongo.QueryBuilder[model.Project]
}

func NewService(db mongo.Database) Service {
	return &service{
		BaseService:         network.NewBaseService(),
		projectQueryBuilder: mongo.NewQueryBuilder[model.Project](db, model.CollectionName),
	}
}

func (s *service) GetProjectSlugExists(slug string) bool {
	filter := bson.M{"slug": slug}
	projection := bson.D{{
		Key:   "created_at",
		Value: 0,
	}}
	opts := options.FindOne().SetProjection(projection)

	_, err := s.projectQueryBuilder.SingleQuery().FindOne(filter, opts)
	return err == nil
}

func (s *service) SaveProject(d *dto.CreateProject) (*model.Project, error) {
	project, err := model.NewProject(d.Title, d.Slug, d.Link, d.ImageURL, d.Type)
	if err != nil {
		return nil, err
	}

	if s.GetProjectSlugExists(d.Slug) {
		return nil, network.NewBadRequestError("slug already exists", err)
	}

	result, err := s.projectQueryBuilder.SingleQuery().InsertAndRetrieveOne(project)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetProjectById(id bson.ObjectID) (*model.Project, error) {
	filter := bson.M{"_id": id}

	project, err := s.projectQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *service) GetProjectBySlug(slug string) (*model.Project, error) {
	filter := bson.M{"slug": slug}

	project, err := s.projectQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// func (s *service) FindPaginatedProjects(p *coredto.Pagination) ([]*model.Project, error) {
// 	filter := bson.M{}

// 	project, err := s.projectQueryBuilder.SingleQuery().FindPaginated(filter, p.Page, p.Limit, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return project, nil
// }

func (s *service) UpdateProject(id bson.ObjectID, d *dto.UpdateProject) error {
	filter := bson.M{"_id": id}
	_, err := s.projectQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return network.NewNotFoundError("project for _id "+id.Hex()+" not found", err)
	}

	update := bson.M{
		"updated_at": time.Now(),
		"title":      d.Title,
		"slug":       d.Slug,
		"link":       d.Link,
		"image_url":  d.ImageURL,
		"type":       d.Type,
	}

	if s.GetProjectSlugExists(d.Slug) {
		return network.NewBadRequestError("slug already exists", err)
	}

	updated := bson.M{"$set": update}

	result, err := s.projectQueryBuilder.SingleQuery().UpdateOne(filter, updated)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return network.NewNotFoundError("project for _id "+id.Hex()+" not found", err)
	}

	return nil
}

func (s *service) DeleteProject(id bson.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := s.projectQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return network.NewNotFoundError("project for _id "+id.Hex()+" not found", err)
	}

	_, err = s.projectQueryBuilder.SingleQuery().DeleteOne(filter)
	if err != nil {
		return err
	}

	return nil
}
