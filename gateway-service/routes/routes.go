package routes

import (
	"github.com/gin-gonic/gin"
	"shopsweb.com/gateway-service/middlewares"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signup", signup)
	server.POST("/login", login)

	//routes that need authentication
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	//authenticated.POST("/events", createEvent)
	//.PUT("/events/:id", updateEvent)
	//authenticated.DELETE("/events/:id", deleteEvent)
	//authenticated.POST("/events/:id/register", registerForEvent)
	//authenticated.DELETE("/events/:id/register", cancelRegistration)
}
