package db

import (
	"context"
	"testing"
	"time"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/stretchr/testify/require"
)

func GenerateRandAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Name:      util.RandomString(5),
		StartDate: time.Now(),
		Email:     util.RandomString(5) + "@gmail.com",
		Password:  util.RandomString(15),
		Weight:    float32(util.RandomInt(150, 250)),
		BodyFat:   float32(util.RandomInt(5, 30)),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.NotNil(t, account.Name)
	require.NotNil(t, account.StartDate)
	require.NotNil(t, account.Email)
	require.NotNil(t, account.Password)
	require.NotNil(t, account.Weight)
	require.NotNil(t, account.BodyFat)
	return account
}

func TestCreateAccount(t *testing.T) {
	GenerateRandAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := GenerateRandAccount(t)
	query, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, account.ID, query.ID)
}

func TestListAccounts(t *testing.T) {
	n := 5
	accounts := make([]Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = GenerateRandAccount(t)
	}

	query, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
		Limit:  int32(n),
		Offset: 0,
	})
	require.NoError(t, err)
	require.Len(t, query, n)
	for _, v := range query {
		require.NotEmpty(t, v)
	}
}

func TestDeleteAccount(t *testing.T) {
	account := GenerateRandAccount(t)

	d, err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}
