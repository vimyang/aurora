package initialize

import (
	"aurora/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRouter() *gin.Engine {
	handler := NewHandle(
		checkProxy(),
		readAccessToken(),
	)

	// 初始化基础前置参数
	handler.InitBasicConfigForChatGPT()

	router := gin.Default()
	// 配置CORS中间件
	router.Use(cors.New(cors.Config{
	    AllowOrigins:     []string{"*"},
	    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	    ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
	    AllowCredentials: true,
	}))

	router.Use(middlewares.Cors)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/auth/session", handler.session)
	router.POST("/auth/refresh", handler.refresh)
	router.OPTIONS("/v1/chat/completions", optionsHandler)

	authGroup := router.Group("").Use(middlewares.Authorization)
	authGroup.POST("/v1/chat/completions", handler.nightmare)
	authGroup.GET("/v1/models", handler.engines)
	authGroup.POST("/backend-api/conversation", handler.chatgptConversation)
	return router
}
