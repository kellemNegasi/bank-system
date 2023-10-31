package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRoutes() {
	// Create a default router
	r := gin.Default()
	// register handlers
	r.POST("/users", server.createUser)
	r.POST("/users/login", server.loginUser)
	r.POST("/accounts", server.createAccount)
	r.GET("/accounts/:id", server.getAccount)
	r.GET("/accounts", server.listAccounts)
	r.POST("/transfers", server.createTransfer)
	server.router = r
}
