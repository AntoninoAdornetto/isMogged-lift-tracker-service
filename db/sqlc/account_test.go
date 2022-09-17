package db

import (
	"context"
	"testing"
	"time"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/stretchr/testify/require"
)

func generateRandAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Lifter:    util.RandomString(6),
		BirthDate: time.Now(),
		Weight:    190,
		StartDate: time.Now(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err) 
	require.NotEmpty(t, acc)

	// todo - Fix Time checks
	require.Equal(t, arg.Lifter, acc.Lifter)
	require.Equal(t, arg.BirthDate.Year(), acc.BirthDate.Year())
	require.Equal(t, arg.Weight, acc.Weight)
	require.Equal(t, arg.StartDate.Year(), acc.StartDate.Year())
	require.NotEmpty(t, acc.ID)

	return acc
}

func TestCreateAccount(t *testing.T) {
	acc := generateRandAccount(t)
	testQueries.DeleteAccount(context.Background(), acc.ID)
}

func TestDeleteAccount(t *testing.T) {
	arg := generateRandAccount(t)
	acc := testQueries.DeleteAccount(context.Background(), arg.ID)
	require.Empty(t, acc)
}

func TestGetAccount(t *testing.T) {
	arg := generateRandAccount(t)
	queryAcc, err := testQueries.GetAccount(context.Background(), arg.ID)
	require.NoError(t, err)
	require.NotEmpty(t, queryAcc)
	testQueries.DeleteAccount(context.Background(), arg.ID)
}

func TestListAccounts(t *testing.T) {
	acc := generateRandAccount(t)
	acc1 := generateRandAccount(t)
	require.NotEmpty(t, acc)
	require.NotEmpty(t, acc1)

	accs, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		Limit:  2,
		Offset: 0,
	})

	require.NoError(t, err)
	require.Len(t, accs, 2)

	testQueries.DeleteAccount(context.Background(), acc.ID)
	testQueries.DeleteAccount(context.Background(), acc1.ID)
}

func TestUpdateAccountWeight(t *testing.T) {
	acc := generateRandAccount(t)
	require.NotEmpty(t, acc)

	testQueries.UpdateAccountWeight(context.Background(), UpdateAccountWeightParams{
		Weight: 195,
		ID:     acc.ID,
	})

	queryAcc, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	require.Equal(t, int(queryAcc.Weight), 195)
	testQueries.DeleteAccount(context.Background(), acc.ID)
}
