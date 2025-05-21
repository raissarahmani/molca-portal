package aws

import (
	"context"

	awsService "github.com/aws/aws-sdk-go/aws"
	awsServiceCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	awsServiceSession "github.com/aws/aws-sdk-go/aws/session"
)

type SessionConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

type Session interface {
	GetInstance() *session
	GetSession() (*awsServiceSession.Session, error)
}

type session struct {
	context context.Context
	config  SessionConfig
}

func NewSession(context context.Context, config SessionConfig) Session {
	return &session{
		context: context,
		config:  config,
	}
}

func (s *session) GetInstance() *session {
	return s
}

func (s *session) GetSession() (*awsServiceSession.Session, error) {
	sess, err := awsServiceSession.NewSession(&awsService.Config{
		Region: awsService.String(s.config.Region),
		Credentials: awsServiceCredentials.NewStaticCredentials(
			s.config.AccessKeyID,
			s.config.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}
