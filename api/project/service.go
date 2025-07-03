package project

import (
	"time"

	"github.com/molca-id/portal-app-api/api/project/dto"
	"github.com/molca-id/portal-app-api/api/project/model"
	"github.com/molca-id/portal-app-api/utils"

	"github.com/molca-id/portal-app-api/arch/mongo"
	"github.com/molca-id/portal-app-api/arch/network"
	"github.com/molca-id/portal-app-api/arch/redis"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Service interface {
	SetProjectDtoCacheById(project *dto.InfoProject) error
	GetProjectDtoCacheById(id bson.ObjectID) (*dto.InfoProject, error)
	SetProjectDtoCacheBySlug(project *dto.InfoProject) error
	GetProjectDtoCacheBySlug(slug string) (*dto.InfoProject, error)
	SaveProject(d *dto.CreateProject) (*model.Project, error)
	GetProjectById(id bson.ObjectID) (*dto.InfoProject, error)
	GetProjectBySlug(slug string) (*dto.InfoProject, error)
	UpdateProject(id bson.ObjectID, d *dto.UpdateProject) error
	GetProjectSlugExists(slug string) bool
	DeleteProject(id bson.ObjectID) error
}

type service struct {
	network.BaseService
	projectQueryBuilder mongo.QueryBuilder[model.Project]
	projectCache        redis.Cache[dto.InfoProject]
}

func NewService(db mongo.Database, store redis.Store) Service {
	return &service{
		BaseService:         network.NewBaseService(),
		projectQueryBuilder: mongo.NewQueryBuilder[model.Project](db, model.CollectionName),
		projectCache:        redis.NewCache[dto.InfoProject](store),
	}
}

func (s *service) SetProjectDtoCacheById(project *dto.InfoProject) error {
	key := "project_" + project.ID.Hex()
	return s.projectCache.SetJSON(key, project, time.Duration(10*time.Minute))
}

func (s *service) SetProjectDtoCacheBySlug(project *dto.InfoProject) error {
	key := "project_" + project.Slug
	return s.projectCache.SetJSON(key, project, time.Duration(10*time.Minute))
}

func (s *service) GetProjectDtoCacheById(id bson.ObjectID) (*dto.InfoProject, error) {
	key := "project_" + id.Hex()
	return s.projectCache.GetJSON(key)
}

func (s *service) GetProjectDtoCacheBySlug(slug string) (*dto.InfoProject, error) {
	key := "project_" + slug
	return s.projectCache.GetJSON(key)
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

func (s *service) GetProjectById(id bson.ObjectID) (*dto.InfoProject, error) {
	filter := bson.M{"_id": id}

	project, err := s.projectQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	result, err := utils.MapTo[dto.InfoProject](project)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetProjectBySlug(slug string) (*dto.InfoProject, error) {
	filter := bson.M{"slug": slug}

	project, err := s.projectQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	result, err := utils.MapTo[dto.InfoProject](project)
	if err != nil {
		return nil, err
	}

	return result, nil
}

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
