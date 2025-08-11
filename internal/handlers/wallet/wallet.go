package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"wallet/pkg/customValidator"
	"wallet/storage"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WalletRequest struct {
	Id            string `json:"id" validate:"required,uuid"`
	OperationType string `json:"operationType" validate:"required,operationType"`
	Amount        int    `json:"amount" validate:"required,gt=0"`
}

type GetWalletRequest struct {
	Id string `validate:"required,uuid"`
}

type WalletRepo interface {
	ChangeAmount(id string, amount int, operation string) error
	Amount(uuid string) (int, error)
}

func ChangeAmountWallet(log *zap.Logger, cv *customValidator.CustomValidator, repo WalletRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		wr := WalletRequest{}

		err := ctx.ShouldBindJSON(&wr)
		if err != nil {
			if errors.Is(err, io.EOF) {
				ctx.JSON(http.StatusBadRequest, "request body is empty")
				return
			}

			var unmarshalTypeErr *json.UnmarshalTypeError
			if errors.As(err, &unmarshalTypeErr) {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("field '%s' must be a %s", unmarshalTypeErr.Field, unmarshalTypeErr.Type),
				})
				return
			}

			ctx.JSON(http.StatusBadRequest, "failed to decode request")
			return
		}

		if err = cv.Validating(&wr); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err = repo.ChangeAmount(wr.Id, wr.Amount, wr.OperationType)
		if err != nil {
			if errors.Is(err, storage.ErrNegativeAmount) {
				ctx.JSON(http.StatusBadRequest, storage.ErrNegativeAmount.Error())
				return
			} else if errors.Is(err, storage.ErrWalletNotFound) {
				ctx.JSON(http.StatusBadRequest, storage.ErrWalletNotFound.Error())
				return
			}
			log.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, "something wrong please retry")
			return
		}
		ctx.JSON(http.StatusOK, "Balance changed")
	}
}

func Amount(log *zap.Logger, cv *customValidator.CustomValidator, repo WalletRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		getWR := GetWalletRequest{Id: id}

		if err := cv.Validating(&getWR); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		amount, err := repo.Amount(getWR.Id)
		if err != nil {
			if err == storage.ErrWalletNotFound {
				ctx.JSON(http.StatusBadRequest, storage.ErrWalletNotFound)
				return
			}
			log.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, storage.ErrWalletNotFound)
			return
		}
		ctx.JSON(http.StatusOK, map[string]int{"Amount :": amount})
	}
}
