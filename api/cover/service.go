package cover

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/molca-id/portal-app-api/api/cover/dto"
	"github.com/molca-id/portal-app-api/arch/aws"
	"github.com/molca-id/portal-app-api/arch/network"
	"github.com/molca-id/portal-app-api/config"
)

type Service interface {
	SaveCover(d *dto.SaveImage) (*dto.ItemImage, error)
}

type service struct {
	network.BaseService
	session aws.Session
	env     *config.Env
}

func NewService(env *config.Env) Service {
	sessConfig := aws.SessionConfig{
		Region:          env.AWSRegion,
		AccessKeyID:     env.AWSAccessKeyID,
		SecretAccessKey: env.AWSSecretAccessKey,
	}
	session := aws.NewSession(context.Background(), sessConfig)
	return &service{
		BaseService: network.NewBaseService(),
		session:     session,
		env:         env,
	}
}

func (s *service) SaveCover(d *dto.SaveImage) (*dto.ItemImage, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	fileExt := path.Ext(d.File.Filename)

	newFilename := timestamp + fileExt

	bucket := aws.NewBucket(s.session)
	err := bucket.PutObject(s.env.AWSBucketName, newFilename, d.File)
	if err != nil {
		return nil, err
	}

	result := &dto.ItemImage{
		ImageUrl: s.env.AWSBucketURL + newFilename,
	}

	return result, nil
}
