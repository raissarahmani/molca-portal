package projects

import (
	"github.com/molca-id/portal-app-api/api/project/model"
	"github.com/molca-id/portal-app-api/api/projects/dto"
	coredto "github.com/molca-id/portal-app-api/arch/dto"
	"github.com/molca-id/portal-app-api/arch/mongo"
	"github.com/molca-id/portal-app-api/arch/network"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Service interface {
	GetPaginatedLatestProjects(fp *dto.FilterPaginatedProject) ([]*dto.ItemProject, error)
	getPublicPaginated(filter bson.M, p *coredto.Pagination) ([]*dto.ItemProject, error)
	getPaginated(filter bson.M, p *coredto.Pagination, opts *options.FindOptionsBuilder) ([]*dto.ItemProject, error)
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

func (s *service) GetPaginatedLatestProjects(fp *dto.FilterPaginatedProject) ([]*dto.ItemProject, error) {
	searchQuery := ""
	if fp.Search != nil {
		searchQuery = *fp.Search
	}

	filter := bson.M{}

	if fp.Type != nil {
		filter["type"] = fp.Type
	}

	filter["$or"] = []bson.M{
		{
			"title": bson.M{
				"$regex":   searchQuery,
				"$options": "i",
			},
		},
		{
			"slug": bson.M{
				"$regex":   searchQuery,
				"$options": "i",
			},
		},
	}

	p := &coredto.Pagination{
		Page:  fp.Page,
		Limit: fp.Limit,
	}

	return s.getPublicPaginated(filter, p)
}

func (s *service) getPublicPaginated(filter bson.M, p *coredto.Pagination) ([]*dto.ItemProject, error) {
	projection := bson.D{{Key: "created_at", Value: 0}}
	opts := options.Find().SetProjection(projection)
	opts.SetSort(bson.D{{Key: "updated_at", Value: -1}})

	return s.getPaginated(filter, p, opts)
}

func (s *service) getPaginated(filter bson.M, p *coredto.Pagination, opts *options.FindOptionsBuilder) ([]*dto.ItemProject, error) {
	projects, err := s.projectQueryBuilder.SingleQuery().FindPaginated(filter, p.Page, p.Limit, opts)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.ItemProject, len(projects))

	for i, p := range projects {
		d, err := dto.NewItemProject(p)
		if err != nil {
			return nil, err
		}

		dtos[i] = d
	}

	return dtos, nil
}
