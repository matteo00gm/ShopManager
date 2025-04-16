package main

import (
	"github.com/gin-gonic/gin"
	"shopsweb.com/gateway-service/routes"
)

func main() {
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8080")
}
