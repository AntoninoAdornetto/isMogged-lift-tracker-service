package db

import (
	"context"
	"testing"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/stretchr/testify/require"
)

func GenerateRandomCategory(t *testing.T) Category {
	var category Category

	category, err := testQueries.CreateCategory(context.Background(), util.RandomString(5))
	require.NoError(t, err)
	require.NotEmpty(t, category)
	require.NotNil(t, category.ID)
	require.NotNil(t, category.Name)
	return category
}

func TestCreateCategory(t *testing.T) {
	GenerateRandomCategory(t)
}

func TestGetCategory(t *testing.T) {
	var category Category

	category, err := testQueries.GetCategory(context.Background(), GenerateRandomCategory(t).ID)
	require.NoError(t, err)
	require.NotEmpty(t, category)
}

func TestListCategories(t *testing.T) {
	n := 5
	categories := make([]Category, n)
	for i := 0; i < n; i++ {
		categories[i] = GenerateRandomCategory(t)
	}

	query, err := testQueries.ListCategories(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(query), n)

	for i := range categories {
		_ = testQueries.DeleteCategory(context.Background(), categories[i].ID)
	}
}

func TestUpdateCategory(t *testing.T) {
	category := GenerateRandomCategory(t)
	newName := util.RandomString(5)

	err := testQueries.UpdateCategory(context.Background(), UpdateCategoryParams{
		ID:   category.ID,
		Name: newName,
	})
	require.NoError(t, err)

	query, err := testQueries.GetCategory(context.Background(), category.ID)
	require.NoError(t, err)
	require.NotEmpty(t, query)
	require.Equal(t, newName, query.Name)
}

func TestDeleteCategory(t *testing.T) {
	category := GenerateRandomCategory(t)

	err := testQueries.DeleteCategory(context.Background(), category.ID)
	require.NoError(t, err)

	query, err := testQueries.GetCategory(context.Background(), category.ID)
	require.Error(t, err)
	require.Empty(t, query)
}
