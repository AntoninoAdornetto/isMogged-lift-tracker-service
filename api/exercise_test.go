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

func TestCreateExercise(t *testing.T) {
	exercise := generateRandExercise()

	testCases := []struct {
		name       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":         exercise.Name,
				"muscle_group": exercise.MuscleGroup,
				"category":     exercise.Category,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateExerciseParams{
					Name:        exercise.Name,
					Category:    exercise.Category,
					MuscleGroup: exercise.MuscleGroup,
				}
				store.EXPECT().CreateExercise(gomock.Any(), gomock.Eq(args)).Times(1).Return(exercise, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateExerciseResponse(t, recorder.Body, exercise)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateExerciseParams{
					Name:        exercise.Name,
					Category:    exercise.Category,
					MuscleGroup: exercise.MuscleGroup,
				}
				store.EXPECT().CreateExercise(gomock.Any(), gomock.Eq(args)).Times(0).Return(db.Exercise{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"name":         exercise.Name,
				"muscle_group": exercise.MuscleGroup,
				"category":     exercise.Category,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateExerciseParams{
					Name:        exercise.Name,
					Category:    exercise.Category,
					MuscleGroup: exercise.MuscleGroup,
				}
				store.EXPECT().CreateExercise(gomock.Any(), gomock.Eq(args)).Times(1).Return(db.Exercise{}, sql.ErrConnDone)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/exercise"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestGetExercise(t *testing.T) {
	exercise := generateRandExercise()

	testCases := []struct {
		name         string
		exerciseName string
		buildStubs   func(store *mockdb.MockStore)
		checkRes     func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			exerciseName: exercise.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetExercise(gomock.Any(), gomock.Eq(exercise.Name)).Times(1).Return(exercise, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateExerciseResponse(t, recorder.Body, exercise)
			},
		},
		{
			name:         "NotFound",
			exerciseName: exercise.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetExercise(gomock.Any(), gomock.Eq(exercise.Name)).Times(1).Return(db.Exercise{}, sql.ErrNoRows)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:         "InternalError",
			exerciseName: exercise.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetExercise(gomock.Any(), gomock.Eq(exercise.Name)).Times(1).Return(db.Exercise{}, sql.ErrConnDone)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/exercise/%s", tc.exerciseName)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestListExercises(t *testing.T) {
	n := 5
	exercises := make([]db.Exercise, n)
	for i := 0; i < n; i++ {
		exercises[i] = generateRandExercise()
	}

	type Query struct {
		PageID   int
		PageSize int
	}

	testCases := []struct {
		name       string
		query      Query
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				PageID:   1,
				PageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.ListExercisesParams{
					Offset: 0,
					Limit:  int32(n),
				}
				store.EXPECT().ListExercises(gomock.Any(), gomock.Eq(args)).Times(1).Return(exercises, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateExercisesResponse(t, recorder.Body, exercises)
			},
		},
		{
			name: "InternalError",
			query: Query{
				PageID:   1,
				PageSize: n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListExercises(gomock.Any(), gomock.Any()).Times(1).Return([]db.Exercise{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				PageID:   -1,
				PageSize: 100000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListExercises(gomock.Any(), gomock.Any()).Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/exercise"
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			qParams := req.URL.Query()
			qParams.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
			qParams.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			req.URL.RawQuery = qParams.Encode()

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestListByMuscleGroup(t *testing.T) {
	exercises, muscleGroup := createMuscleGroupExercises()

	type Query struct {
		PageSize    int
		PageID      int
		MuscleGroup string
	}

	testCases := []struct {
		name       string
		query      Query
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				PageID:      1,
				PageSize:    5,
				MuscleGroup: muscleGroup,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.ListByMuscleGroupParams{
					MuscleGroup: muscleGroup,
					Limit:       5,
					Offset:      0,
				}
				store.EXPECT().ListByMuscleGroup(gomock.Any(), gomock.Eq(args)).Times(1).Return(exercises, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateExercisesResponse(t, recorder.Body, exercises)
			},
		},
		{
			name: "InternalError",
			query: Query{
				PageID:      1,
				PageSize:    5,
				MuscleGroup: muscleGroup,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListByMuscleGroup(gomock.Any(), gomock.Any()).Times(1).Return([]db.Exercise{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				PageID:      -10000,
				PageSize:    5,
				MuscleGroup: muscleGroup,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListByMuscleGroup(gomock.Any(), gomock.Any()).Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/exercise/group/%s", tc.query.MuscleGroup)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			qParams := req.URL.Query()
			qParams.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
			qParams.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			req.URL.RawQuery = qParams.Encode()

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestUpdateExercise(t *testing.T) {
	srcExercise := generateRandExercise()
	shiftExercise := db.Exercise{
		Name:        util.RandomString(5),
		Category:    util.RandomString(5),
		MuscleGroup: util.RandomString(5),
		ID:          srcExercise.ID,
	}

	testCases := []struct {
		name       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":         shiftExercise.Name,
				"muscle_group": shiftExercise.MuscleGroup,
				"category":     shiftExercise.Category,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.UpdateExerciseParams{
					Name:    srcExercise.Name,
					Column1: shiftExercise.Name,
					Column2: shiftExercise.MuscleGroup,
					Column3: shiftExercise.Category,
				}
				store.EXPECT().UpdateExercise(gomock.Any(), gomock.Eq(args)).Return(shiftExercise, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateExerciseResponse(t, recorder.Body, shiftExercise)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"name":         shiftExercise.Name,
				"muscle_group": shiftExercise.MuscleGroup,
				"category":     shiftExercise.Category,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.UpdateExerciseParams{
					Name:    srcExercise.Name,
					Column1: shiftExercise.Name,
					Column2: shiftExercise.MuscleGroup,
					Column3: shiftExercise.Category,
				}
				store.EXPECT().UpdateExercise(gomock.Any(), gomock.Eq(args)).Times(1).Return(db.Exercise{}, sql.ErrConnDone)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/exercise/%s", srcExercise.Name)
			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestDeleteExercise(t *testing.T) {
	exercise := generateRandExercise()

	testCases := []struct {
		name         string
		exerciseName string
		buildStubs   func(store *mockdb.MockStore)
		checkRes     func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			exerciseName: exercise.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteExercise(gomock.Any(), gomock.Eq(exercise.Name)).Times(1).Return(nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name:         "InternalError",
			exerciseName: exercise.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteExercise(gomock.Any(), gomock.Eq(exercise.Name)).Times(1).Return(sql.ErrConnDone)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/exercise/%s", tc.exerciseName)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func updateExercise(src *db.Exercise, shift db.Exercise) (*db.Exercise, db.Exercise) {
	var orignalVals db.Exercise
	orignalVals = *src
	*src = shift
	return src, orignalVals
}

func generateRandExercise() db.Exercise {
	return db.Exercise{
		Name:        util.RandomString(5),
		Category:    util.RandomString(5),
		ID:          int32(util.RandomInt(1, 200)),
		MuscleGroup: util.RandomString(5),
	}
}

func createMuscleGroupExercises() ([]db.Exercise, string) {
	n := 5
	mg := util.RandomString(n)
	exercises := make([]db.Exercise, n)
	for i := 0; i < n; i++ {
		exercises[i] = db.Exercise{
			Name:        util.RandomString(n),
			MuscleGroup: mg,
		}
	}
	return exercises, mg
}

func validateExerciseResponse(t *testing.T, body *bytes.Buffer, exercise db.Exercise) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resExercise db.Exercise
	err = json.Unmarshal(data, &resExercise)
	require.NoError(t, err)
	require.Equal(t, exercise, resExercise)
}

func validateExercisesResponse(t *testing.T, body *bytes.Buffer, exercise []db.Exercise) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resExercises []db.Exercise
	err = json.Unmarshal(data, &resExercises)
	require.NoError(t, err)
	require.Equal(t, exercise, resExercises)
}
