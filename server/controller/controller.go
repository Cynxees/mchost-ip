package controller

import (
	"context"
	"mchost-ip/server/api"
	"net/http"
	"time"

	pb "mchost-ip/server/pb"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupHandlers(router *gin.Engine, server *api.Server) {

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.LoadHTMLFiles("index.html")

	router.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.POST("/api/get-instance", func(c *gin.Context) {
		getIpFleetHandler(c, server)
	})

	router.POST("/api/create", func(c *gin.Context) {
		createIpFleetHandler(c, server)
	})

	router.POST("/api/use", func(c *gin.Context) {
		useIpHandler(c, server)
	})

	router.POST("/api/unuse", func(c *gin.Context) {
		unuseIpHandler(c, server)
	})

	router.POST("/api/delete", func(c *gin.Context) {
		deleteIpHandler(c, server)
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

// CreateIpFleetIp godoc
// @Summary Create Ip Fleet Ip
// @Tags Ip
// @Accept json
// @Produce json
// @Param			requestBody	body		pb.CreateIpRequest	true	"Request Body"
// @Success 200 {string} Helloworld
// @Router /create [post]
func createIpFleetHandler(c *gin.Context, server *api.Server) {
	ctx := context.Background()

	var request pb.CreateIpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := server.CreateIp(ctx, &request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}

// GetIpFleetIp godoc
// @Summary Get Ip Fleet Instances
// @Tags Ip
// @Accept json
// @Produce json
// @Param			requestBody	body		pb.GetIpRequest	true	"Request Body"
// @Success 200 {string} Helloworld
// @Router /get-instance [post]
func getIpFleetHandler(c *gin.Context, server *api.Server) {
	ctx := context.Background()

	var request pb.GetIpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := server.GetIp(ctx, &request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}

// UseIp godoc
// @Summary Use Ip
// @Tags Ip
// @Accept json
// @Produce json
// @Param			requestBody	body		pb.UseIpRequest	true	"Request Body"
// @Success 200 {string} Helloworld
// @Router /use [post]
func useIpHandler(c *gin.Context, server *api.Server) {
	ctx := context.Background()

	var request pb.UseIpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := server.UseIp(ctx, &request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}

// UnuseIp godoc
// @Summary Unuse Ip
// @Tags Ip
// @Accept json
// @Produce json
// @Param			requestBody	body		pb.UnuseIpRequest	true	"Request Body"
// @Success 200 {string} Helloworld
// @Router /unuse [post]
func unuseIpHandler(c *gin.Context, server *api.Server) {
	ctx := context.Background()

	var request pb.UnuseIpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := server.UnuseIp(ctx, &request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}

// DeleteIp godoc
// @Summary Delete Ip
// @Tags Ip
// @Accept json
// @Produce json
// @Param			requestBody	body		pb.DeleteIpRequest	true	"Request Body"
// @Success 200 {string} Helloworld
// @Router /delete [post]
func deleteIpHandler(c *gin.Context, server *api.Server) {
	ctx := context.Background()

	var request pb.DeleteIpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := server.DeleteIp(ctx, &request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}