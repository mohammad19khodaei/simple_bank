package middlewares_test

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mohammad19khodaei/simple_bank/utils"
)

var (
	config utils.Config
)

func TestMain(m *testing.M) {
	cfg, err := utils.LoadConfig("../..", "app.testing")
	if err != nil {
		log.Fatal("could not load config", err)
	}
	config = cfg

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
