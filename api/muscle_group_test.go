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

	mockdb "github.com/AntoninoAdornetto/lift_tracker/db/mock"
	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateMuscleGroup(t *testing.T) {
	muscleGroup := generateRandMuscleGroup()

	testCases := []struct {
		name       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name": muscleGroup.Name,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateMuscleGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Return(muscleGroup, nil)
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
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateMuscleGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Return(db.MuscleGroup{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/muscle_group"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)

		})
	}
}

func TestGetMuscleGroup(t *testing.T) {
	muscleGroup := generateRandMuscleGroup()

	testCases := []struct {
		name        string
		muscleGroup string
		buildStubs  func(store *mockdb.MockStore)
		checkRes    func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			muscleGroup: muscleGroup.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMuscleGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Return(muscleGroup, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:        "InternalError",
			muscleGroup: muscleGroup.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMuscleGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Return(db.MuscleGroup{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/muscle_group/%s", tc.muscleGroup)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestListMuscleGroups(t *testing.T) {
	muscleGroups := generateRandMuscleGroups()

	testCases := []struct {
		name       string
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMuscleGroups(gomock.Any()).Return(muscleGroups, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateMuscleGroupsResponse(t, recorder.Body, muscleGroups)
			},
		},
		{
			name: "InternalError",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMuscleGroups(gomock.Any()).Return([]db.MuscleGroup{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/muscle_group"
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestDeleteMuscleGroup(t *testing.T) {
	muscleGroup := generateRandMuscleGroup()

	testCases := []struct {
		name        string
		muscleGroup string
		buildStubs  func(store *mockdb.MockStore)
		checkRes    func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			muscleGroup: muscleGroup.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteGroup(gomock.Any(), gomock.Eq(muscleGroup.Name)).Return(muscleGroup, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateMuscleGroupResponse(t, recorder.Body, muscleGroup)
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/muscle_group/%s", tc.muscleGroup)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

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
