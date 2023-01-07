package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AntoninoAdornetto/isMogged-lift-tracker-service/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func addAuthHeader(
	t *testing.T,
	request *http.Request,
	tokenCreator token.Maker,
	authorizationType string,
	userId uuid.UUID,
	duration time.Duration) {
	token, err := tokenCreator.CreateToken(userId, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthenticationMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		checkRes      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, "unsupported", uuid.New(), time.Minute)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, "", uuid.New(), time.Minute)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), -time.Minute)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			route := "/auth"
			server.router.GET(route, authenticationMiddleware(server.tokenCreator), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, route, nil)
			require.NoError(t, err)

			tc.configureAuth(t, request, server.tokenCreator)
			server.router.ServeHTTP(recorder, request)
			tc.checkRes(t, recorder)
		})
	}
}
