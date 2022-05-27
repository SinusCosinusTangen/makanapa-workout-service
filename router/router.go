package router

import (
	"github.com/gin-gonic/gin"

	"workout-webservice/controller/api"
	"workout-webservice/controller/middleware"
)

func InitializeRouter() (router *gin.Engine) {
	router = gin.Default()

	apiRoute := router.Group("/api")
	apiRoute.Use(
		middleware.CorsMiddleware,
	)
	{
		workout := apiRoute.Group("/workout")
		{
			workout.POST("/fetch-workout-by-time", api.FetchWorkoutByTime)
			workout.POST("/fetch-workout-history", api.FetchWorkoutHistory)
			workout.GET("/db-seed", api.SeedDB)
		}
	}

	return router
}
