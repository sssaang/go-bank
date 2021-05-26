package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	db "github.com/sssaang/simplebank/db/sqlc"
	testdb "github.com/sssaang/simplebank/db/test"
	"github.com/sssaang/simplebank/db/util"
	"github.com/sssaang/simplebank/token"
	"github.com/stretchr/testify/require"
)

func TestMakeTransfer(t *testing.T){
	amount := int64(10)
	user1, _ := randomUser(t)
	account1 := randomAccount(user1.Username)
	user2, _ := randomUser(t)
	account2 := randomAccount(user2.Username)
	account2.Currency = account1.Currency
	user3, _ := randomUser(t)
	account3 := randomAccount(user3.Username)
	switch account1.Currency {
	case util.USD, util.KRW: 
		account3.Currency = util.EUR
	case util.EUR:
		account3.Currency = util.USD
	}

	testCases := []struct {
		name string
		body gin.H
		setupAuth func(t *testing.T, request *http.Request, tokenManager token.TokenManager)
		buildStubs func(store *testdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Make a successful transfer",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": account1.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager){
				addAuthorization(t, request, tokenManager, AUTHORIZATION_TYPE_BEARER, user1.Username, time.Minute)
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
				Times(1).Return(account1, nil)
				
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
				Times(1).Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID: account2.ID,
					Amount: 10,
				}

				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Eq(arg)).
				Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Currency Mismatch",
			body: gin.H {
				"from_account_id": account1.ID,
				"to_account_id": account3.ID,
				"amount": amount,
				"currency": account1.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager){
				addAuthorization(t, request, tokenManager, AUTHORIZATION_TYPE_BEARER, user1.Username, time.Minute)
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
				Times(1).Return(account1, nil)
				
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account3.ID)).
				Times(1).Return(account3, nil)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Currency",
			body: gin.H {
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": "Invalid Currency",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager){
				addAuthorization(t, request, tokenManager, AUTHORIZATION_TYPE_BEARER, user1.Username, time.Minute)
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
				Times(0)
				
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
				Times(0)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          -amount,
				"currency":        account1.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager){
				addAuthorization(t, request, tokenManager, AUTHORIZATION_TYPE_BEARER, user1.Username, time.Minute)
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetAccountError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        account1.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager){
				addAuthorization(t, request, tokenManager, AUTHORIZATION_TYPE_BEARER, user1.Username, time.Minute)
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Any()).
				Times(1).Return(db.Account{}, sql.ErrConnDone)

				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Any()).
				Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TransferTxError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        account1.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager){
				addAuthorization(t, request, tokenManager, AUTHORIZATION_TYPE_BEARER, user1.Username, time.Minute)
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
				Times(1).Return(account1, nil)
				
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
				Times(1).Return(account2, nil)
				
				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Any()).
				Times(1).Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized user transfer",
			body: gin.H {
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": account1.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager){
				addAuthorization(t, request, tokenManager, AUTHORIZATION_TYPE_BEARER, "asdf", time.Minute)
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Any()).
				Times(1).Return(account1, nil)
				
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
				Times(0)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "No Authorization",
			body: gin.H {
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": account1.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenManager token.TokenManager){
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Any()).
				Times(0)
				
				store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
				Times(0)

				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/transfer", bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenManager)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}