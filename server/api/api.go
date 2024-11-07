package api

import (
	"mchost-spot-instance/server/config"
	manager "mchost-spot-instance/server/jwt"
	"mchost-spot-instance/server/pb"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server struct {
	pb.SpotServiceServer
	Db      *gorm.DB
	Logger  *logrus.Logger
	Manager *manager.JWTManager
	Config  *config.Config
}
