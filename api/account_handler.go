package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR IRR"`
}

func (s *server) createAccountHandler(ctx *gin.Context) {
	var request CreateAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, s.errorResponse(err))
	}

	account, err := s.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    request.Owner,
		Currency: request.Currency,
		Balance:  0,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
	}

	ctx.JSON(http.StatusCreated, account)
}
