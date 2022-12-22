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
		Email:     util.RandomEmail(),
		Password:  util.RandomString(15),
		Weight:    float32(util.RandomInt(150, 250)),
		BodyFat:   float32(util.RandomInt(5, 30)),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Name, args.Name)
	require.NotZero(t, account.StartDate)
	require.Equal(t, account.Email, args.Email)
	require.Equal(t, account.Password, args.Password)
	require.Equal(t, account.Weight, args.Weight)
	require.Equal(t, account.BodyFat, args.BodyFat)
	require.True(t, account.PasswordChangedAt.IsZero())
	return account
}

func TestCreateAccount(t *testing.T) {
	GenerateRandAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := GenerateRandAccount(t)
	query, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, query)

	require.Equal(t, account.ID, query.ID)
	require.Equal(t, account.BodyFat, query.BodyFat)
	require.Equal(t, account.Weight, query.Weight)
	require.Equal(t, account.Password, query.Password)
	require.Equal(t, account.Email, query.Email)
	require.Equal(t, account.Name, query.Name)
	require.WithinDuration(t, account.StartDate, query.StartDate, time.Second)
	require.WithinDuration(t, account.PasswordChangedAt, query.PasswordChangedAt, time.Second)
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
