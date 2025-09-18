package config

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	// server
	GoMode     string `mapstructure:"GO_MODE"`
	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort uint16 `mapstructure:"SERVER_PORT"`
	// database
	DBHost         string `mapstructure:"DB_HOST"`
	DBName         string `mapstructure:"DB_NAME"`
	DBPort         uint16 `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBUserPwd      string `mapstructure:"DB_USER_PWD"`
	DBMinPoolSize  uint16 `mapstructure:"DB_MIN_POOL_SIZE"`
	DBMaxPoolSize  uint16 `mapstructure:"DB_MAX_POOL_SIZE"`
	DBQueryTimeout uint16 `mapstructure:"DB_QUERY_TIMEOUT_SEC"`
	// redis
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort uint16 `mapstructure:"REDIS_PORT"`
	RedisPwd  string `mapstructure:"REDIS_PASSWORD"`
	RedisDB   int    `mapstructure:"REDIS_DB"`
	// AWS
	AWSAccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AWSRegion          string `mapstructure:"AWS_REGION"`
	AWSBucketName      string `mapstructure:"AWS_BUCKET_NAME"`
	AWSBucketURL       string `mapstructure:"AWS_BUCKET_URL"`
	// keys
	RSAPrivateKeyPath string `mapstructure:"RSA_PRIVATE_KEY_PATH"`
	RSAPublicKeyPath  string `mapstructure:"RSA_PUBLIC_KEY_PATH"`
	// Token
	AccessTokenValiditySec  uint64 `mapstructure:"ACCESS_TOKEN_VALIDITY_SEC"`
	RefreshTokenValiditySec uint64 `mapstructure:"REFRESH_TOKEN_VALIDITY_SEC"`
	TokenIssuer             string `mapstructure:"TOKEN_ISSUER"`
	TokenAudience           string `mapstructure:"TOKEN_AUDIENCE"`
	// JWKs
	JWKsURL string `mapstructure:"JWKS_URL"`
	// analytics
	GoogleApplicationCredentials string `mapstructure:"GOOGLE_APPLICATION_CREDENTIALS"`
	GAPropertyID                 string `mapstructure:"GA_PROPERTY_ID"`
}

func NewEnv(filename string, override bool) *Env {
	env := Env{}
	viper.SetConfigFile(filename)

	if override {
		viper.AutomaticEnv()
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading environment file", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Error loading environment file", err)
	}

	return &env
}
