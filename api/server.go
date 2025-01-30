package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
)

type server struct {
	store  *db.Store
	server *gin.Engine
}

func NewServer(pool *pgxpool.Pool) *server {
	server := &server{
		store: db.NewStore(pool),
	}
	r := gin.Default()
	r.POST("/accounts", server.createAccountHandler)
	r.GET("/accounts/:id", server.getAccountHandler)

	server.server = r
	return server
}

func (s *server) Start(address string) error {
	return s.server.Run(address)
}

func (s *server) errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
