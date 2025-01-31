package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
)

type server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *server {
	server := &server{store: store}

	r := gin.Default()
	r.POST("/accounts", server.createAccountHandler)
	r.GET("/accounts/:id", server.getAccountHandler)
	r.GET("/accounts", server.ListAccountsHandler)
	r.POST("/transfer", server.transferHandler)

	server.router = r
	return server
}

func (s *server) Start(address string) error {
	return s.router.Run(address)
}

func (s *server) GetRouter() *gin.Engine {
	return s.router
}

func (s *server) errorResponse(err error) gin.H {
	return gin.H{
		"message": err.Error(),
	}
}
