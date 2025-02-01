package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *server) createAccountHandler(ctx *gin.Context) {
	var request createAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, s.errorResponse(err))
		return
	}

	account, err := s.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    request.Owner,
		Currency: request.Currency,
		Balance:  0,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				ctx.JSON(http.StatusConflict, s.errorResponse(fmt.Errorf("Account already exists with the same owner and currency")))
			case pgerrcode.ForeignKeyViolation:
				ctx.JSON(http.StatusBadRequest, s.errorResponse(fmt.Errorf("Owner does not exist")))
			}
			return
		}
		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountParams struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *server) getAccountHandler(ctx *gin.Context) {
	var params getAccountParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, s.errorResponse(err))
		return
	}

	account, err := s.store.GetAccount(ctx, params.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, s.errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsParams struct {
	Page    int32 `form:"page" binding:"required,min=1"`
	PerPage int32 `form:"per_page" binding:"required,min=5,max=10"`
}

func (s *server) ListAccountsHandler(ctx *gin.Context) {
	var params listAccountsParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, s.errorResponse(err))
		return
	}

	accounts, err := s.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  params.PerPage,
		Offset: (params.Page - 1) * params.PerPage,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
