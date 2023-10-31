package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/kellemNegasi/bank-system/db/mock"
	db "github.com/kellemNegasi/bank-system/db/sqlc"
	token "github.com/kellemNegasi/bank-system/token/pasto"
	"github.com/kellemNegasi/bank-system/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountId     int64
		buildStub     func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker)
		checkResponse func(t *testing.T, recorder httptest.ResponseRecorder)
	}{
		{
			name:      "ok",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuth(t, tokenMaker, request, time.Minute, user.Username, authorizationTypeBearer)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireAccountBodyMath(t, recorder.Body, account)
			},
		},

		{
			name:      "NotFound",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuth(t, tokenMaker, request, time.Minute, user.Username, authorizationTypeBearer)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},

			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		{
			name:      "Internal Error",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuth(t, tokenMaker, request, time.Minute, user.Username, authorizationTypeBearer)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name:      "Invalid ID",
			accountId: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuth(t, tokenMaker, request, time.Minute, user.Username, authorizationTypeBearer)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name:      "Unauthorized",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuth(t, tokenMaker, request, time.Minute, "unauthorized_user", authorizationTypeBearer)
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			accountId: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		// TODO more test cases to come.
	}

	for _, tc := range testCases {
		// initialize
		ctrl := gomock.NewController(t)
		store := mockdb.NewMockStore(ctrl)

		// build stub
		tc.buildStub(store)

		// start test server
		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()

		url := fmt.Sprintf("/accounts/%d", tc.accountId)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		tc.setupAuth(t, request, *server.Maker)

		server.router.ServeHTTP(recorder, request)
		tc.checkResponse(t, *recorder)

	}

}
func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       int64(util.RandInt(1, 1000)),
		Owner:    owner,
		Currency: util.RandCurrency(),
		Balance:  util.RandomMoney(20, 200),
	}
}

func requireAccountBodyMath(t *testing.T, body *bytes.Buffer, expectedAccount db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var actualAccount db.Account
	err = json.Unmarshal(data, &actualAccount)
	require.NoError(t, err)
	require.Equal(t, expectedAccount, actualAccount)
}
