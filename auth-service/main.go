package main

import (
	"github.com/gin-gonic/gin"
	"shopsweb.com/auth-service/db"
	"shopsweb.com/auth-service/routes"
)

func main() {
	db.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8080")
}
