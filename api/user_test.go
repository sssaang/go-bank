package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	db "github.com/sssaang/simplebank/db/sqlc"
	testdb "github.com/sssaang/simplebank/db/test"
	"github.com/sssaang/simplebank/db/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name string
		body gin.H
		buildStubs func(store *testdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Create an user",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"full_name": user.FullName,
				"email": user.Email,
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Times(1).
				Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "Invalid Username",
			body: gin.H{
				"username": "invalid",
				"password": password,
				"full_name": user.FullName,
				"email": user.Email,
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Email",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"full_name": user.FullName,
				"email": "invalidemail",
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Too short password",
			body: gin.H{
				"username": user.Username,
				"password": "1234",
				"full_name": user.FullName,
				"email": user.Email,
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Create a duplicate user",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"full_name": user.FullName,
				"email": user.Email,
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Times(1).
				Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "Internal Error",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"full_name": user.FullName,
				"email": user.Email,
			},
			buildStubs: func(store *testdb.MockStore) {
				store.EXPECT().
				CreateUser(gomock.Any(), gomock.Any()).
				Times(1).
				Return(db.User{}, sql.ErrConnDone)
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	username := util.RandomEmail()
	password = util.RandomString(10)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User {
		Username: username,
		HashedPassword: hashedPassword,
		FullName: util.RandomOwner(),
		Email: username,
	}

	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}