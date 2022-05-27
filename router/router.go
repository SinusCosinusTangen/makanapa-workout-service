package router

import (
	"github.com/gin-gonic/gin"

	"workout-webservice/controller/api"
)

func InitializeRouter() (router *gin.Engine) {
	router = gin.Default()

	// router.Use(
	// 	middleware.CorsMiddleware,
	// )

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	apiRoute := router.Group("/api")
	{
		workout := apiRoute.Group("/workout")
		{
			workout.POST("/fetch-workout-by-time", api.FetchWorkoutByTime)
			workout.GET("/fetch-workout-history", api.FetchWorkoutHistory)
		}
	}

	return router
}
