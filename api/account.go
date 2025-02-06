package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mohammad19khodaei/simple_bank/api/middlewares"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *server) createAccountHandler(ctx *gin.Context) {
	owner := ctx.MustGet(middlewares.AuthUsernameKey).(string)

	var request createAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, s.errorResponse(err))
		return
	}

	account, err := s.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    owner,
		Currency: request.Currency,
		Balance:  0,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				ctx.JSON(http.StatusConflict, s.errorResponse(errors.New("Account already exists with the same owner and currency")))
			case pgerrcode.ForeignKeyViolation:
				ctx.JSON(http.StatusBadRequest, s.errorResponse(errors.New("Owner does not exist")))
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
	username := ctx.MustGet(middlewares.AuthUsernameKey).(string)

	if account.Owner != username {
		ctx.JSON(http.StatusForbidden, s.errorResponse(errors.New("forbidden: account does not belong to you")))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsParams struct {
	Page    int32 `form:"page" binding:"omitempty,min=1"`
	PerPage int32 `form:"per_page" binding:"omitempty,min=5,max=10"`
}

func (s *server) ListAccountsHandler(ctx *gin.Context) {
	var params listAccountsParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, s.errorResponse(err))
		return
	}

	page := int32(1)
	if params.Page != 0 {
		page = params.Page
	}

	perPage := int32(10)
	if params.PerPage != 0 {
		perPage = params.PerPage
	}

	username := ctx.MustGet(middlewares.AuthUsernameKey).(string)

	accounts, err := s.store.ListAccounts(ctx, db.ListAccountsParams{
		Owner:  username,
		Limit:  perPage,
		Offset: (page - 1) * perPage,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
