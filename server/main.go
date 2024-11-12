package main

import (
	"fmt"
	"mchost-ip/server/api"
	awsManager "mchost-ip/server/aws"
	"mchost-ip/server/config"
	controller "mchost-ip/server/controller"
	jwtManager "mchost-ip/server/jwt"

	// "mchost-ip/server/lib/rabbitmq"
	"mchost-ip/server/models"
	"mchost-ip/server/pb"
	"net"

	"mchost-ip/www/docs"

	elasticLog "gopkg.in/sohlich/elogrus.v7"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	fmt.Println("starting server")
	appConfig := config.InitConfig(".env")

	docs.SwaggerInfo.Title = "Ip Instance API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", appConfig.CurrentAddress, appConfig.AppPort)
	docs.SwaggerInfo.BasePath = "/api"

	esLogger := logrus.New()
	esLogger.SetFormatter(&ecslogrus.Formatter{})
	esLogger.SetLevel(logrus.InfoLevel)

	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false))
	if err != nil {
		esLogger.Fatalf("Failed to create elasticsearch client: %v", err)
	}

	hook, err := elasticLog.NewAsyncElasticHook(client, appConfig.AppName, logrus.InfoLevel, "go-auth-logs")
	if err != nil {
		logrus.Fatalf("Failed to create Elasticsearch hook: %v", err)
	}
	esLogger.AddHook(hook)

	dsn := "user:password@tcp(127.0.0.1:" + appConfig.DbPort + ")/mchost_ip?charset=utf8mb4&parseTime=True&loc=Local"

	esLogger.Info(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		esLogger.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&models.Ip{})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", appConfig.MicroservicePort))
	if err != nil {
		esLogger.Fatalf("failed to listen: %v", err)
	}

	// rabbitmq.InitRabbitMq()

	// // Start the RabbitMQ consumer in a goroutine
	// go rabbitmq.StartConsumer("orders")

	grpcServer := grpc.NewServer()
	server := &api.Server{
		Db:         db,
		Logger:     esLogger,
		JWTManager: jwtManager.NewJWTManager(appConfig.AppKey, 3600, esLogger),
		AppConfig:  appConfig,
		AWSManager: awsManager.NewAWSManager(appConfig.AwsAccessKeyId, appConfig.AwsAccessKeySecret),
	}

	router := gin.Default()

	controller.SetupHandlers(router, server)

	go router.Run(fmt.Sprintf(":%s", appConfig.AppPort))
	pb.RegisterIpServiceServer(grpcServer, server)

	if err := grpcServer.Serve(lis); err != nil {
		esLogger.Fatalf("failed to serve: %v", err)
	}
}
