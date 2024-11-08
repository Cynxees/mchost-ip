package api

import (
	aws "mchost-spot-instance/server/aws"
	"mchost-spot-instance/server/config"
	manager "mchost-spot-instance/server/jwt"
	"mchost-spot-instance/server/pb"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server struct {
	pb.SpotServiceServer
	Db         *gorm.DB
	Logger     *logrus.Logger
	JWTManager *manager.JWTManager
	AppConfig  *config.Config
	AWSManager *aws.AWSManager
}
