package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kellemNegasi/bank-system/db/sqlc"
	token "github.com/kellemNegasi/bank-system/token/pasto"
	"github.com/kellemNegasi/bank-system/util"
)

// Server represents the HTTP server that serves client requests.
type Server struct {
	config util.Config
	store  db.Store
	Maker  *token.PasetoMaker
	router *gin.Engine
}

// New returns a new Server object.
func New(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPastoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		config: config,
		Maker:  tokenMaker,
		store:  store,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRoutes()
	return server, err
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
