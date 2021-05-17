package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/sssaang/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// Creates a new HTTP server, setup routing and return the server
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfer", server.createTransfer)

	server.router = router
	return server
}

// Starts to run the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H { // gin.H is a shortcut for map[string] interface
	return gin.H {"error": err.Error()}
}