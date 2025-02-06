package middlewares_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mohammad19khodaei/simple_bank/api/middlewares"
	"github.com/mohammad19khodaei/simple_bank/token"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setAuthHeader func(t *testing.T, maker token.Maker, req *http.Request)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "without authorization header",
			setAuthHeader: func(t *testing.T, _ token.Maker, _ *http.Request) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "without bearer keyword",
			setAuthHeader: func(t *testing.T, _ token.Maker, req *http.Request) {
				req.Header.Set("Authorization", "invalid-token")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "with invalid header",
			setAuthHeader: func(t *testing.T, tokenMaker token.Maker, req *http.Request) {
				token, err := tokenMaker.GenerateToken("username", config.TokenDuration)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", "unsupported", token))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "with expired token",
			setAuthHeader: func(t *testing.T, tokenMaker token.Maker, req *http.Request) {
				token, err := tokenMaker.GenerateToken("username", -config.TokenDuration)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, token))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ok",
			setAuthHeader: func(t *testing.T, tokenMaker token.Maker, req *http.Request) {
				token, err := tokenMaker.GenerateToken("username", config.TokenDuration)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, token))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	tokenMaker, err := token.NewPasetoMaker(config.SecretKey)
	require.NoError(t, err)
	r := gin.Default()
	authUrl := "/auth"
	r.GET(authUrl, middlewares.AuthMiddleware(tokenMaker), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, authUrl, nil)

			tc.setAuthHeader(t, tokenMaker, request)
			r.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
