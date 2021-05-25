package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sssaang/simplebank/token"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenManager token.TokenManager,
	authorization_type string,
	username string,
	duration time.Duration,
) {
	token, err := tokenManager.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	authorizationHeader := fmt.Sprintf("%s %s", AUTHORIZATION_TYPE_BEARER, token)
	request.Header.Set(AUTHORIZATION_HEADER, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name string
		setupAuth func(t *testing.T, request *http.Request, tokenManager token.TokenManager)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := NewTestServer(t, nil)

			server.router.GET(
				"/auth",
				authMiddleware(server.tokenManager),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "auth", nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenManager)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}