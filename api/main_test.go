package api_test

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
	"github.com/mohammad19khodaei/simple_bank/utils"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func createRandomAccount(currency ...string) db.Account {
	acc := db.Account{
		ID:      int32(utils.RandomInt(1, 1000)),
		Owner:   utils.RandomOwner(),
		Balance: utils.RandomMoney(),
	}

	// Use provided currency if specified, otherwise random
	if len(currency) > 0 {
		acc.Currency = currency[0]
	} else {
		acc.Currency = utils.RandomCurrency()
	}

	return acc
}
