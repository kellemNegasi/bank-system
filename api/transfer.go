package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/kellemNegasi/bank-system/db/sqlc"
	token "github.com/kellemNegasi/bank-system/token/pasto"
	"github.com/shopspring/decimal"
)

type createTransferRequest struct {
	FromAccountID int64            `json:"from_account_id" binding:"required"`
	ToAccountID   int64            `json:"to_account_id" binding:"required"`
	Currency      string           `json:"currency" binding:"required,currency"`
	Amount        *decimal.Decimal `json:"amount" binding:"required"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check for self transfer
	if req.FromAccountID == req.ToAccountID {
		err := fmt.Errorf("self transfer is prohibited! from account %d to account %d", req.FromAccountID, req.ToAccountID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.isValidAccount(ctx, req.FromAccountID, req.Currency, true, req.Amount)

	if !valid {
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.PasetoPayload)
	if fromAccount.Owner != payload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.isValidAccount(ctx, req.ToAccountID, req.Currency, false, nil)
	if !valid {
		return
	}

	args := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount.String(),
	}

	result, err := server.store.TransferTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) isValidAccount(ctx *gin.Context, accountID int64, currency string, checkBalance bool, amount *decimal.Decimal) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err = fmt.Errorf("currency mismatch. currency %s vs  currency %s. \n account id %d", account.Currency, currency, accountID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	if checkBalance && amount != nil {
		// Check if the account has enough balance for the transfer
		Balance, err := decimal.NewFromString(account.Balance)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return account, false
		}

		if Balance.LessThan(*amount) {
			err = fmt.Errorf("insufficient balance in account %d for transfer", accountID)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return account, false
		}
	}

	return account, true
}
