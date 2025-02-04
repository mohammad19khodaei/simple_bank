package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mohammad19khodaei/simple_bank/api/validators"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
	"github.com/mohammad19khodaei/simple_bank/token"
	"github.com/mohammad19khodaei/simple_bank/utils"
)

type server struct {
	store      db.Store
	tokenMaker token.Maker
	config     utils.Config
	router     *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.SecretKey)
	if err != nil {
		return nil, err
	}
	server := &server{
		tokenMaker: tokenMaker,
		config:     config,
		store:      store,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validators.CurrencyValidator)
	}

	server.registerRouter()

	return server, nil
}

func (s *server) registerRouter() {
	r := gin.Default()

	r.POST("/users", s.createUserHandler)
	r.POST("/users/login", s.login)

	r.POST("/accounts", s.createAccountHandler)
	r.GET("/accounts/:id", s.getAccountHandler)
	r.GET("/accounts", s.ListAccountsHandler)
	r.POST("/transfer", s.transferHandler)

	s.router = r
}

func (s *server) Start(address string) error {
	return s.router.Run(address)
}

func (s *server) Router() *gin.Engine {
	return s.router
}

func (s *server) errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
