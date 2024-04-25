package db

import (
	"context"
	"testing"

	"github.com/pdadu/learn/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  int64(util.RandomAmount()),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotEmpty(t, account.ID)
	require.NotEmpty(t, account.CreatedAt)
	return account
}
func Test_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func Test_GetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
	require.Equal(t, account1.Owner, account2.Owner)
}

func Test_UpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: int64(util.RandomAmount()),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
	require.Equal(t, account1.Owner, account2.Owner)
}

func Test_DeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

}

func Test_ListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccount(context.Background(), arg)
	require.Len(t, accounts, 5)
	require.NoError(t, err)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
