package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/stretchr/testify/require"
)

func CreateTestAccount() Account {
	arg := CreateAccountParams{
		Lifter:    util.RandomString(6),
		BirthDate: time.Now(),
		Weight:    190,
		StartDate: time.Now(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), arg)

	if err != nil {
		os.Exit(1)
	}

	return acc
}

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Lifter:    util.RandomString(6),
		BirthDate: time.Now(),
		Weight:    190,
		StartDate: time.Now(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err) // checks for nil error
	require.NotEmpty(t, account)

	// todo - Fix Time checks
	require.Equal(t, arg.Lifter, account.Lifter)
	require.Equal(t, arg.BirthDate.Year(), account.BirthDate.Year())
	require.Equal(t, arg.Weight, account.Weight)
	require.Equal(t, arg.StartDate.Year(), account.StartDate.Year())
	require.NotEmpty(t, account.ID)
}

func TestDeleteAccount(t *testing.T) {
	arg := CreateTestAccount()
	acc := testQueries.DeleteAccount(context.Background(), arg.ID)
	require.Empty(t, acc)
}

func TestGetAccount(t *testing.T) {
	arg := CreateTestAccount()
	queryAcc, err := testQueries.GetAccount(context.Background(), arg.ID)
	require.NoError(t, err)
	require.NotEmpty(t, queryAcc)
	testQueries.DeleteAccount(context.Background(), arg.ID)
}

func TestListAccounts(t *testing.T) {
	acc := CreateTestAccount()
	acc1 := CreateTestAccount()
	require.NotEmpty(t, acc)
	require.NotEmpty(t, acc1)

	accs, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		Limit:  2,
		Offset: 0,
	})

	require.NoError(t, err)
	require.Len(t, accs, 2)
}

func TestUpdateAccountWeight(t *testing.T) {
	acc := CreateTestAccount()
	require.NotEmpty(t, acc)

	testQueries.UpdateAccountWeight(context.Background(), UpdateAccountWeightParams{
		Weight: 195,
		ID:     acc.ID,
	})

	queryAcc, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	require.Equal(t, int(queryAcc.Weight), 195)
}
