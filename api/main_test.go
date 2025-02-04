package api_test

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
	"github.com/mohammad19khodaei/simple_bank/utils"
)

var (
	config utils.Config
)

func TestMain(m *testing.M) {
	cfg, err := utils.LoadConfig("../", "app.testing")
	if err != nil {
		log.Fatal("could not load config", err)
	}
	config = cfg

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

func createRandomUser(password string) db.User {
	hashedPassword, _ := utils.HashPassword(password)
	return db.User{
		Username:          utils.RandomOwner(),
		HashedPassword:    hashedPassword,
		FullName:          utils.RandomOwner(),
		Email:             utils.RandomEmail(),
		PasswordChangedAt: pgtype.Timestamptz{},
		CreatedAt:         pgtype.Timestamptz{},
	}
}
