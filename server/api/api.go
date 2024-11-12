package api

import (
	aws "mchost-ip/server/aws"
	"mchost-ip/server/config"
	manager "mchost-ip/server/jwt"
	"mchost-ip/server/pb"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server struct {
	pb.IpServiceServer
	Db         *gorm.DB
	Logger     *logrus.Logger
	JWTManager *manager.JWTManager
	AppConfig  *config.Config
	AWSManager *aws.AWSManager
	Redis      *redis.Client
}
