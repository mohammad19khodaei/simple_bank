package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mohammad19khodaei/simple_bank/api/validators"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
)

type server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *server {
	server := &server{store: store}

	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validators.CurrencyValidator)
	}

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
		"error": err.Error(),
	}
}
