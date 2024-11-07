package main

import (
	"context"
	"fmt"
	"mchost-spot-instance/server/api"
	"mchost-spot-instance/server/config"
	manager "mchost-spot-instance/server/jwt"

	// "mchost-spot-instance/server/lib/rabbitmq"
	"mchost-spot-instance/server/models"
	"mchost-spot-instance/server/pb"
	"net"

	"mchost-spot-instance/www/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// elasticLog "gopkg.in/sohlich/elogrus.v7"

	"github.com/gin-gonic/gin"
	// "github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	appConfig := config.InitConfig(".env")

	docs.SwaggerInfo.Title = "Go Auth API"
	docs.SwaggerInfo.Description = "API documentation for Go Auth Service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", appConfig.AppPort)
	docs.SwaggerInfo.BasePath = "/api"

	esLogger := logrus.New()
	esLogger.SetFormatter(&ecslogrus.Formatter{})
	esLogger.SetLevel(logrus.InfoLevel)

	// client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false))
	// if err != nil {
	// 	esLogger.Fatalf("Failed to create elasticsearch client: %v", err)
	// }

	// hook, err := elasticLog.NewAsyncElasticHook(client, serviceName, logrus.InfoLevel, "go-auth-logs")
	// if err != nil {
	// 	logrus.Fatalf("Failed to create Elasticsearch hook: %v", err)
	// }
	// esLogger.AddHook(hook)

	dsn := "user:password@tcp(127.0.0.1:" + appConfig.DbPort + ")/mchost_spot_instance?charset=utf8mb4&parseTime=True&loc=Local"

	esLogger.Info(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		esLogger.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&models.SpotInstance{})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", appConfig.MicroservicePort))
	if err != nil {
		esLogger.Fatalf("failed to listen: %v", err)
	}

	// rabbitmq.InitRabbitMq()

	// // Start the RabbitMQ consumer in a goroutine
	// go rabbitmq.StartConsumer("orders")

	grpcServer := grpc.NewServer()
	service := &api.Server{
		Db:      db,
		Logger:  esLogger,
		Manager: manager.NewJWTManager(appConfig.AppKey, 3600, esLogger),
		Config:  appConfig,
	}

	router := gin.Default()

	// Route to serve the Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/register", func(c *gin.Context) {
		ctx := context.Background()

		resp, err := service.RegisterSpotInstance(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Return the response
		c.JSON(200, resp)
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	go router.Run(fmt.Sprintf(":%s", appConfig.AppPort))
	pb.RegisterSpotServiceServer(grpcServer, service)

	if err := grpcServer.Serve(lis); err != nil {
		esLogger.Fatalf("failed to serve: %v", err)
	}
}
