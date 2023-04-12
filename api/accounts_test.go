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
	mockdb "github.com/julysNICK/simplebank/db/mock"
	db "github.com/julysNICK/simplebank/db/sqlc"
	"github.com/julysNICK/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestListAccounts(t *testing.T) {
	user := randomUser()
	user2 := randomUser()
	user3 := randomUser()
	lists := []db.Account{
		randomAccount(user.Username),
		randomAccount(user2.Username),
		randomAccount(user3.Username),
	}

	testCases := []struct {
		name          string
		limit         int32
		offset        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			limit:  5,
			offset: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(1).Return(lists, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccountList(t, recorder.Body, lists)
			},
		},
		{
			name:   "INVALID_LIMIT",
			limit:  0,
			offset: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		}, {
			name:   "INTERNAL_ERROR",
			limit:  10,
			offset: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(1).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			store := mockdb.NewMockStore(gomock.NewController(t))
			testCase.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts?page_id=%d&limit=%d", testCase.limit, testCase.offset), nil)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}
func TestCreateAccountAPI(t *testing.T) {
	user := randomUser()
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		body          io.Reader
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: mustNewJSONBody(t, createAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
			}),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "Invalid JSON",
			body: bytes.NewBufferString("invalid json"),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "INTERNAL SERVER ERROR",
			body: mustNewJSONBody(t, createAccountRequest{}),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := mockdb.NewMockStore(gomock.NewController(t))
			server := newTestServer(t,store)

			tc.buildStubs(store)

			recorder := httptest.NewRecorder()
			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, tc.body)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
func TestGetAccountAPI(t *testing.T) {
	user := randomUser()
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		AccountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			AccountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NOT FOUND",
			AccountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

			},
		},
		{
			name:      "INTERNAL SERVER ERROR",
			AccountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name:      "INVALID ACCOUNT ID",
			AccountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
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

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t,store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.AccountID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func randomAccount(ownerName string) db.Account {
	return db.Account{
		ID:       utils.RandomINt(1, 1000),
		Owner:    ownerName,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account

	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account, gotAccount)
}

func mustNewJSONBody(t *testing.T, data interface{}) *bytes.Buffer {
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(data)
	require.NoError(t, err)
	return body
}

func requireBodyMatchAccountList(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account

	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)

	require.Equal(t, accounts, gotAccounts)
}
