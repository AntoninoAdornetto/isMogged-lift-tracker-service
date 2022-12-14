package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/AntoninoAdornetto/isMogged-lift-tracker-service/db/mock"
	db "github.com/AntoninoAdornetto/isMogged-lift-tracker-service/db/sqlc"
	"github.com/AntoninoAdornetto/isMogged-lift-tracker-service/token"
	"github.com/AntoninoAdornetto/isMogged-lift-tracker-service/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory(t *testing.T) {
	category := generateRandCategory()

	testCases := []struct {
		name          string
		body          gin.H
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name": category.Name,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(category.Name)).Times(1).Return(category, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateCategoryResponse(t, recorder.Body, category)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(category.Name)).Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"name": category.Name,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(category.Name)).Times(1).Return(db.Category{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"name": category.Name,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(category.Name)).Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/category"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestGetCategory(t *testing.T) {
	category := generateRandCategory()

	testCases := []struct {
		name          string
		categoryId    int16
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			categoryId: category.ID,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCategory(gomock.Any(), gomock.Eq(category.ID)).Times(1).Return(category, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateCategoryResponse(t, recorder.Body, category)
			},
		},
		{
			name:       "InternalError",
			categoryId: category.ID,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCategory(gomock.Any(), gomock.Eq(category.ID)).Times(1).Return(db.Category{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "Unauthorized",
			categoryId: category.ID,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCategory(gomock.Any(), gomock.Eq(category.ID)).Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/category/%d", tc.categoryId)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestListCategories(t *testing.T) {
	n := 5
	categories := make([]db.Category, n)
	for i := 0; i < n; i++ {
		categories[i] = generateRandCategory()
	}

	testCases := []struct {
		name          string
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListCategories(gomock.Any()).Times(1).Return(categories, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateCategoriesResponse(t, recorder.Body, categories)
			},
		},
		{
			name: "InternalError",
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListCategories(gomock.Any()).Times(1).Return([]db.Category{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListCategories(gomock.Any()).Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/category"
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestUpdateCategory(t *testing.T) {
	category := generateRandCategory()

	testCases := []struct {
		name          string
		body          gin.H
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name": category.Name,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.UpdateCategoryParams{
					Name: category.Name,
					ID:   category.ID,
				}
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Eq(args)).Times(1).Return(nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Any()).Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"name": category.Name,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.UpdateCategoryParams{
					Name: category.Name,
					ID:   category.ID,
				}
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Eq(args)).Times(1).Return(sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"name": category.Name,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.UpdateCategoryParams{
					Name: category.Name,
					ID:   category.ID,
				}
				store.EXPECT().UpdateCategory(gomock.Any(), gomock.Eq(args)).Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/category/%d", category.ID)
			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestDeleteCategory(t *testing.T) {
	category := generateRandCategory()

	testCases := []struct {
		name          string
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCategory(gomock.Any(), gomock.Eq(category.ID)).Times(1).Return(nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},

		{
			name: "Unauthorized",
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteCategory(gomock.Any(), gomock.Eq(category.ID)).Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/category/%d", category.ID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func generateRandCategory() db.Category {
	return db.Category{
		Name: util.RandomString(5),
		ID:   int16(util.RandomInt(1, 100)),
	}
}

func validateCategoryResponse(t *testing.T, body *bytes.Buffer, category db.Category) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resCategory db.Category
	err = json.Unmarshal(data, &resCategory)
	require.NoError(t, err)
	require.Equal(t, category, resCategory)
}

func validateCategoriesResponse(t *testing.T, body *bytes.Buffer, category []db.Category) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resCategorys []db.Category
	err = json.Unmarshal(data, &resCategorys)
	require.NoError(t, err)
	require.Equal(t, category, resCategorys)
}
