package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRoutes() {
	// Create a default router
	r := gin.Default()
	// register handlers
	r.POST("/users", server.createUser)
	r.POST("/users/login", server.loginUser)

	// Authorized endpoints
	routGroups := r.Group("/").Use(authMiddleware(*server.Maker))
	routGroups.POST("/accounts", server.createAccount)
	routGroups.GET("/accounts/:id", server.getAccount)
	routGroups.GET("/accounts", server.listAccounts)
	routGroups.POST("/transfers", server.createTransfer)

	server.router = r
}
