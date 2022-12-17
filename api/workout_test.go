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
	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateWorkout(t *testing.T) {
	workout := generateRandWorkout()

	testCases := []struct {
		name       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"start_time": workout.StartTime.UnixMilli(),
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateWorkoutParams{
					UserID:    workout.UserID,
					StartTime: util.FormatMSEpoch(workout.StartTime.UnixMilli()),
				}
				store.EXPECT().CreateWorkout(gomock.Any(), gomock.Eq(args)).Times(1).Return(workout, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateWorkout(gomock.Any(), gomock.Any()).Times(0).Return(db.Workout{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"start_time": workout.StartTime.UnixMilli(),
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateWorkoutParams{
					UserID:    workout.UserID,
					StartTime: util.FormatMSEpoch(workout.StartTime.UnixMilli()),
				}
				store.EXPECT().CreateWorkout(gomock.Any(), gomock.Eq(args)).Times(1).Return(db.Workout{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/workout/%s", workout.UserID)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)

		})
	}
}

func TestGetWorkout(t *testing.T) {
	workouts := generateRandWorkouts()

	testCases := []struct {
		name       string
		workoutID  uuid.UUID
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			workoutID: workouts[0].ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetWorkout(gomock.Any(), gomock.Any()).Times(1).Return(workouts, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			workoutID: workouts[0].ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetWorkout(gomock.Any(), gomock.Any()).Times(1).Return([]db.GetWorkoutRow{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/workout/%s", workouts[0].ID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)

		})
	}
}

func TestListWorkouts(t *testing.T) {
	n := 5
	workouts := make([]db.Workout, n)
	for i := 0; i < n; i++ {
		workouts[i] = generateRandWorkout()
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
				args := db.ListWorkoutsParams{
					Limit:  int32(n),
					Offset: 0,
					UserID: workouts[0].UserID,
				}
				store.EXPECT().ListWorkouts(gomock.Any(), gomock.Eq(args)).Times(1).Return(workouts, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
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

			url := fmt.Sprintf("/workout/history/%s", workouts[0].UserID)
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

func TestUpdateDurationEnd(t *testing.T) {
	workout := generateRandWorkout()

	testCases := []struct {
		name       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"finish_time": workout.FinishTime.UnixMilli(),
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.UpdateFinishTimeParams{
					FinishTime: util.FormatMSEpoch(time.Now().UnixMilli()),
					ID:         workout.ID,
				}
				store.EXPECT().UpdateFinishTime(gomock.Any(), gomock.Eq(args)).Times(1).Return(workout, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
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

			url := fmt.Sprintf("/workout/%s", workout.ID)
			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestDeleteWorkout(t *testing.T) {
	workout := generateRandWorkout()

	testCases := []struct {
		name       string
		workoutID  uuid.UUID
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "NoContent",
			workoutID: workout.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteWorkout(gomock.Any(), gomock.Eq(workout.ID)).Times(1).Return(nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
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

			url := fmt.Sprintf("/workout/%s", tc.workoutID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)

		})
	}
}

func generateRandWorkout() db.Workout {
	return db.Workout{
		ID:         uuid.New(),
		StartTime:  util.FormatMSEpoch(time.Now().UnixMilli()),
		FinishTime: util.FormatMSEpoch(time.Now().UnixMilli()),
		UserID:     uuid.New(),
	}
}

func generateRandWorkouts() []db.GetWorkoutRow {
	n := 5
	workouts := make([]db.GetWorkoutRow, n)
	for i := 0; i < n; i++ {
		workouts[i] = db.GetWorkoutRow{
			ID:         uuid.New(),
			UserID:     uuid.New(),
			FinishTime: time.Now(),
			StartTime:  time.Now(),
		}
	}
	return workouts
}

func validateWorkoutResponse(t *testing.T, body *bytes.Buffer, lift db.Workout) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resWorkout db.Workout
	err = json.Unmarshal(data, &resWorkout)
	require.NoError(t, err)
	require.Equal(t, lift, resWorkout)
}

func validateWorkoutsResponse(t *testing.T, body *bytes.Buffer, lift []db.Workout) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resWorkouts []db.Workout
	err = json.Unmarshal(data, &resWorkouts)
	require.NoError(t, err)
	require.Equal(t, lift, resWorkouts)
}
