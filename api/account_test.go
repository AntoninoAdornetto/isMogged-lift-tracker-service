package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/AntoninoAdornetto/lift_tracker/db/mock"
	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/AntoninoAdornetto/lift_tracker/token"
	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	userID := uuid.New()
	account := generateRandAccount(userID)

	testCases := []struct {
		name          string
		accountID     uuid.UUID
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, userID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateAccountResponse(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Unauthorized",
			accountID: account.ID,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(0)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%s", tc.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(t, recorder)
		})
	}
}

func TestListAccounts(t *testing.T) {
	var lastAccount db.Account
	userID := uuid.New()
	n := 5
	accounts := make([]db.Account, n)
	for i := 0; i < n; i++ {
		lastAccount = generateRandAccount(userID)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, lastAccount.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.ListAccountsParams{
					ID:     lastAccount.ID,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(args)).Times(1).Return(accounts, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateAccountsResponse(t, recorder.Body, accounts)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, lastAccount.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(1).Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, lastAccount.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   -1,
				pageSize: 10000,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, lastAccount.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			query: Query{
				pageID:   -1,
				pageSize: 10000,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/accounts"
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			qParams := req.URL.Query()
			qParams.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			qParams.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			req.URL.RawQuery = qParams.Encode()

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func generateRandAccount(userID uuid.UUID) db.Account {
	password := util.RandomString(10)
	hashedPassword, _ := util.HashPassword(password)

	return db.Account{
		ID:       userID,
		Name:     util.RandomString(10),
		Email:    util.RandomString(5) + "@gmail.com",
		Password: hashedPassword,
		Weight:   float32(util.RandomInt(150, 220)),
		BodyFat:  float32(util.RandomInt(8, 25)),
	}
}

func validateAccountResponse(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resAccount db.Account
	err = json.Unmarshal(data, &resAccount)
	require.NoError(t, err)
	require.Equal(t, account.ID, resAccount.ID)
	require.Equal(t, account.Name, resAccount.Name)
	require.Equal(t, account.Email, resAccount.Email)
	require.Equal(t, account.Weight, resAccount.Weight)
	require.Equal(t, account.BodyFat, resAccount.BodyFat)
	require.WithinDuration(t, account.StartDate, resAccount.StartDate, time.Second)
}

func validateAccountsResponse(t *testing.T, body *bytes.Buffer, account []db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resAccounts []db.Account
	err = json.Unmarshal(data, &resAccounts)
	require.NoError(t, err)
	require.Equal(t, account, resAccounts)
}
