package controller

import (
	"context"
	"mchost-spot-instance/server/api"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
  pb "mchost-spot-instance/server/pb"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupHandlers(router *gin.Engine, server *api.Server) {

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/api/launch", func(c *gin.Context) {
		launchSpotFleetHandler(c, server)
	})

	router.POST("/api/get-instance", func(c *gin.Context) {
		getSpotFleetHandler(c, server)
	})

	router.GET("/api/ping", Helloworld)
	router.POST("/api/pong", Pong)
}

func Pong(g *gin.Context) {

	currentTime := time.Now()
	_, err := g.GetRawData()
	time.Sleep(1 * time.Second)
	duration := time.Since(currentTime)

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"time": duration.Seconds()})
}


// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /ping [get]
func Helloworld(g *gin.Context) {
	time.Sleep(1 * time.Second)
	g.JSON(http.StatusOK, "helloworld")
}

// LaunchSpotFleet godoc
// @Summary Launch Spot Fleet Instances
// @Schemes
// @Tags Spot
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /launch [get]
func launchSpotFleetHandler(c *gin.Context, server *api.Server) {
	ctx := context.Background()
	
	resp, err := server.LaunchSpotFleet(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}

// GetSpotFleet godoc
// @Summary Get Spot Fleet Instances
// @Tags Spot
// @Accept json
// @Produce json
// @Param			requestBody	body		pb.GetSpotFleetRequest	true	"Request Body"
// @Success 200 {string} Helloworld
// @Router /get-instance [post]
func getSpotFleetHandler(c *gin.Context, server *api.Server) {
	ctx := context.Background()

	var request pb.GetTemplateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := server.GetSpotFleet(ctx, request.SpotInstanceId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}