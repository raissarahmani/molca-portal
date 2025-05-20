package project

import (
	"github.com/molca-id/portal-app-api/api/project/dto"
	"github.com/molca-id/portal-app-api/api/project/model"
	coredto "github.com/molca-id/portal-app-api/arch/dto"
	"github.com/molca-id/portal-app-api/arch/mongo"
	"github.com/molca-id/portal-app-api/arch/network"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service interface {
	SaveProject(d *dto.CreateProject) (*model.Project, error)
	GetProjectById(id bson.ObjectID) (*model.Project, error)
	GetProjectBySlug(slug string) (*model.Project, error)
	FindPaginatedProjects(p *coredto.Pagination) ([]*model.Project, error)
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

func (s *service) SaveProject(d *dto.CreateProject) (*model.Project, error) {
	project, err := model.NewProject(d.Title, d.Slug, d.Link, d.ImageURL, d.Type)
	if err != nil {
		return nil, err
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

func (s *service) FindPaginatedProjects(p *coredto.Pagination) ([]*model.Project, error) {
	filter := bson.M{}

	project, err := s.projectQueryBuilder.SingleQuery().FindPaginated(filter, p.Page, p.Limit, nil)
	if err != nil {
		return nil, err
	}

	return project, nil
}
