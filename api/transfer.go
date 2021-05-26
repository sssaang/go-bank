package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/sssaang/simplebank/db/sqlc"
	"github.com/sssaang/simplebank/token"
)


type transferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID int64 `json:"to_account_id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required,min=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) makeTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check if the the account has the same currency

	fromAccount, isValid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !isValid {
		return
	}

	authPayload := ctx.MustGet(AUTHORIZATION_PAYLOAD).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("the user has no access to the from account")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, isValid = server.validAccount(ctx, req.ToAccountID, req.Currency) 
	if !isValid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return db.Account{}, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return db.Account{}, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: the currency of the account is %s while the currency of the transfer is %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return db.Account{}, false
	}

	return account, true
}