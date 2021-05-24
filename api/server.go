package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/sssaang/simplebank/db/sqlc"
	"github.com/sssaang/simplebank/db/util"
	"github.com/sssaang/simplebank/token"
)

type Server struct {
	config util.Config
	store  db.Store
	tokenManager token.TokenManager
	router *gin.Engine
}

// NewServer creates a new HTTP server, setup routing and return the server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenManager, err := token.NewPasetoManager(config.PasetoSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token manager %w", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenManager: tokenManager,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/user", server.createUser)
	router.GET("/user/:username", server.getUser)
	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfer", server.createTransfer)

	server.router = router
	return server, nil
}

// Starts to run the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H { // gin.H is a shortcut for map[string] interface
	return gin.H {"error": err.Error()}
}