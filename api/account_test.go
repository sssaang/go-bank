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

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	db "github.com/sssaang/simplebank/db/sqlc"
	testdb "github.com/sssaang/simplebank/db/test"
	"github.com/sssaang/simplebank/db/util"
	"github.com/sssaang/simplebank/token"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name string
		body gin.H
		buildStubs func(store *testdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Create an account",
			body: gin.H{
				"owner": account.Owner,
				"currency": account.Currency,
				"balance": account.Balance,
			},
			buildStubs: func(store *testdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner: account.Owner,
					Currency: account.Currency,
				}
				store.EXPECT().
				CreateAccount(gomock.Any(), gomock.Eq(arg)).
				Times(1).
				Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "Non-existent user",
			body: gin.H {
				"owner": "non-existent user",
				"currency": account.Currency,
				"balance": account.Balance,
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				CreateAccount(gomock.Any(), gomock.Any()).
				Times(1).
				Return(db.Account{}, &pq.Error{Code: "23503"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "Invalid Currency",
			body: gin.H {
				"owner": account.Owner,
				"currency": "invalid currency",
				"balance": account.Balance,
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				CreateAccount(gomock.Any(), gomock.Any()).
				Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			body: gin.H{
				"owner": account.Owner,
				"currency": account.Currency,
				"balance": account.Balance,
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				CreateAccount(gomock.Any(), gomock.Any()).
				Times(1).
				Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := testdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/account", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}


func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name string
		accountID int64
		setupAuth func(t *testing.T, request *http.Request, tokenManager token.TokenManager)
		buildStubs func(store *testdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	} {
			{
				name: "Get an existing account",
				accountID: account.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager){
					addAuthorization(t, request, tokenManager, AUTHORIZATION_TYPE_BEARER, "test_user", time.Minute)
				},
				buildStubs: func(store *testdb.MockStore) {
					store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
					requireBodyMatchAccount(t, recorder.Body, account)
				},
			},
			{
				name: "Get an account that does not exist",
				accountID: account.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager){
					addAuthorization(t, request, tokenManager, AUTHORIZATION_TYPE_BEARER, "test_user", time.Minute)
				},
				buildStubs: func(store *testdb.MockStore) {
					store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusNotFound, recorder.Code)
				},
			},
			{
				name: "Connection Error",
				accountID: account.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager){
					addAuthorization(t, request, tokenManager, AUTHORIZATION_TYPE_BEARER, "test_user", time.Minute)
				},
				buildStubs: func(store *testdb.MockStore) {
					store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
			},
			{
				name: "Invalid ID",
				accountID: -12,
				setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager){
					addAuthorization(t, request, tokenManager, AUTHORIZATION_TYPE_BEARER, "test_user", time.Minute)
				},
				buildStubs: func(store *testdb.MockStore) {
					store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
	}

	for i := range testCases {
		tc := testCases[i]
		
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			
			store := testdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/account/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenManager)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount(username string) db.Account {
	return db.Account {
		ID: util.RandomInt(1, 10000),
		Owner: username,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}