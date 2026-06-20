package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/naki/dispatch-service/controllers"
	"github.com/naki/dispatch-service/transport/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "dispatch"})
	})

	api := r.Group("/api/v1")
	api.Use(middlewares.AuthMiddleware())
	{
		nurse := api.Group("/availability")
		nurse.Use(middlewares.RequireRole("nurse"))
		{
			nurse.POST("/online", controllers.GoOnline)
			nurse.POST("/offline", controllers.GoOffline)
			nurse.PUT("/location", controllers.UpdateLocation)
		}

		dispatch := api.Group("/dispatch")
		{
			dispatch.GET("/nurses", middlewares.RequireRole("nurse", "super_admin"), controllers.GetAvailableNurses)
			dispatch.POST("/match", middlewares.RequireRole("super_admin"), controllers.ManualDispatch)
			dispatch.GET("/history/:booking_id", middlewares.RequireRole("nurse", "customer", "super_admin"), controllers.GetDispatchHistory)
			dispatch.GET("/recent", middlewares.RequireRole("super_admin"), controllers.GetRecentDispatches)
		}
	}

	return r
}
