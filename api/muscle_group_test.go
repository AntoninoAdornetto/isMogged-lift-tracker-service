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

	mockdb "github.com/AntoninoAdornetto/lift_tracker/db/mock"
	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/AntoninoAdornetto/lift_tracker/token"
	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateMuscleGroup(t *testing.T) {
	muscleGroup := generateRandMuscleGroup()

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
				"name": muscleGroup.Name,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateMuscleGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Times(1).Return(muscleGroup, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"name": muscleGroup.Name,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateMuscleGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Times(1).Return(db.MuscleGroup{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"name": muscleGroup.Name,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateMuscleGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Times(0)
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

			url := "/muscle_group"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)

		})
	}
}

func TestGetMuscleGroup(t *testing.T) {
	muscleGroup := generateRandMuscleGroup()

	testCases := []struct {
		name          string
		muscleGroup   string
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			muscleGroup: muscleGroup.Name,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMuscleGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Times(1).Return(muscleGroup, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:        "InternalError",
			muscleGroup: muscleGroup.Name,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMuscleGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Times(1).Return(db.MuscleGroup{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:        "Unauthorized",
			muscleGroup: muscleGroup.Name,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMuscleGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Times(0)
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

			url := fmt.Sprintf("/muscle_group/%s", tc.muscleGroup)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestListMuscleGroups(t *testing.T) {
	muscleGroups := generateRandMuscleGroups()

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
				store.EXPECT().GetMuscleGroups(gomock.Any()).Times(1).Return(muscleGroups, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateMuscleGroupsResponse(t, recorder.Body, muscleGroups)
			},
		},
		{
			name: "InternalError",
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMuscleGroups(gomock.Any()).Times(1).Return([]db.MuscleGroup{}, sql.ErrConnDone)
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
				store.EXPECT().GetMuscleGroups(gomock.Any()).Times(0)
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

			url := "/muscle_group"
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestDeleteMuscleGroup(t *testing.T) {
	muscleGroup := generateRandMuscleGroup()

	testCases := []struct {
		name          string
		muscleGroup   string
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			muscleGroup: muscleGroup.Name,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Times(1).Return(muscleGroup, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateMuscleGroupResponse(t, recorder.Body, muscleGroup)
			},
		},
		{
			name:        "Unauthorized",
			muscleGroup: muscleGroup.Name,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Times(0)
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

			url := fmt.Sprintf("/muscle_group/%s", tc.muscleGroup)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func generateRandMuscleGroup() db.MuscleGroup {
	return db.MuscleGroup{
		Name: util.RandomString(5),
		ID:   int16(util.RandomInt(1, 150)),
	}
}

func generateRandMuscleGroups() []db.MuscleGroup {
	n := 5
	muscleGroups := make([]db.MuscleGroup, n)
	for i := 0; i < n; i++ {
		muscleGroups[i] = db.MuscleGroup{
			Name: util.RandomString(5),
			ID:   int16(util.RandomInt(1, 150)),
		}
	}
	return muscleGroups
}

func validateMuscleGroupResponse(t *testing.T, body *bytes.Buffer, lift db.MuscleGroup) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resMuscleGroup db.MuscleGroup
	err = json.Unmarshal(data, &resMuscleGroup)
	require.NoError(t, err)
	require.Equal(t, lift, resMuscleGroup)
}

func validateMuscleGroupsResponse(t *testing.T, body *bytes.Buffer, lift []db.MuscleGroup) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resMuscleGroups []db.MuscleGroup
	err = json.Unmarshal(data, &resMuscleGroups)
	require.NoError(t, err)
	require.Equal(t, lift, resMuscleGroups)
}
