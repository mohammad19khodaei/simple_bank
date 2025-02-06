package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/mohammad19khodaei/simple_bank/api/middlewares"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
)

type transferRequest struct {
	FromAccountID int32 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int32 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
}

func (s *server) transferHandler(ctx *gin.Context) {
	var request transferRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, s.errorResponse(err))
		return
	}

	fromAccount, err := s.store.GetAccount(ctx, request.FromAccountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, s.errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
		return
	}

	username := ctx.MustGet(middlewares.AuthUsernameKey).(string)
	if fromAccount.Owner != username {
		ctx.JSON(http.StatusForbidden, s.errorResponse(errors.New("forbidden: account does not belong to you")))
		return
	}

	if fromAccount.Balance < request.Amount {
		ctx.JSON(http.StatusPaymentRequired, s.errorResponse(errors.New("insufficient balance")))
		return
	}

	toAccount, err := s.store.GetAccount(ctx, request.ToAccountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, s.errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
		return
	}

	if fromAccount.Currency != toAccount.Currency {
		ctx.JSON(http.StatusBadRequest, s.errorResponse(errors.New(fmt.Sprintf("from account currency %s mismatch to account currency %s", fromAccount.Currency, toAccount.Currency))))
		return
	}

	transfer, err := s.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        request.Amount,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, transfer)
}
