package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kellemNegasi/bank-system/db/sqlc"
)

// Server represents the HTTP server that serves client requests.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// New returns a new Server object.
func New(store db.Store) *Server {
	server := &Server{
		store: store,
	}

	// Create a default router
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// register handlers
	r.POST("/users", server.createUser)
	r.POST("/accounts", server.createAccount)
	r.GET("/accounts/:id", server.getAccount)
	r.GET("/accounts", server.listAccounts)
	r.POST("/transfers", server.createTransfer)

	server.router = r

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
