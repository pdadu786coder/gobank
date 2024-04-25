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

	"github.com/golang/mock/gomock"
	mockdb "github.com/pdadu/learn/db/mock"
	db "github.com/pdadu/learn/db/sqlc"
	"github.com/pdadu/learn/db/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatch(t, *recorder.Body, account)
			},
		},
		{
			name:      "noID",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
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
			name:      "internalservererror",
			accountID: account.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		// {
		// 	name:      "invalidinput",
		// 	accountID: 0,
		// 	buildStub: func(store *mockdb.MockStore) {
		// 		store.EXPECT().
		// 			GetAccount(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			crtl := gomock.NewController(t)
			defer crtl.Finish()
			store := mockdb.NewMockStore(crtl)
			tc.buildStub(store)
			server := NewServer(store)
			recoder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recoder, request)
			tc.checkResponse(t, recoder)
		})

	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  int64(util.RandomAmount()),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatch(t *testing.T, body bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(&body)
	require.NoError(t, err)
	var gotaccount db.Account
	err = json.Unmarshal(data, &gotaccount)
	require.NoError(t, err)
	require.Equal(t, gotaccount, account)

}
