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

func TestCreateLift(t *testing.T) {
	lift := generateRandLift()

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
				"exercise_name": lift.ExerciseName,
				"weight":        lift.WeightLifted,
				"reps":          lift.Reps,
				"user_id":       lift.UserID,
				"workout_id":    lift.WorkoutID,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateLiftParams{
					ExerciseName: lift.ExerciseName,
					WeightLifted: lift.WeightLifted,
					Reps:         lift.Reps,
					UserID:       lift.UserID,
					WorkoutID:    lift.WorkoutID,
				}
				store.EXPECT().CreateLift(gomock.Any(), gomock.Eq(args)).Times(1).Return(lift, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateLiftResponse(t, recorder.Body, lift)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateLiftParams{
					ExerciseName: lift.ExerciseName,
					WeightLifted: lift.WeightLifted,
					Reps:         lift.Reps,
					UserID:       lift.UserID,
					WorkoutID:    lift.WorkoutID,
				}
				store.EXPECT().CreateLift(gomock.Any(), gomock.Eq(args)).Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"exercise_name": lift.ExerciseName,
				"weight":        lift.WeightLifted,
				"reps":          lift.Reps,
				"user_id":       lift.UserID,
				"workout_id":    lift.WorkoutID,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateLiftParams{
					ExerciseName: lift.ExerciseName,
					WeightLifted: lift.WeightLifted,
					Reps:         lift.Reps,
					UserID:       lift.UserID,
					WorkoutID:    lift.WorkoutID,
				}
				store.EXPECT().CreateLift(gomock.Any(), gomock.Eq(args)).Times(1).Return(db.Lift{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"exercise_name": lift.ExerciseName,
				"weight":        lift.WeightLifted,
				"reps":          lift.Reps,
				"user_id":       lift.UserID,
				"workout_id":    lift.WorkoutID,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateLiftParams{
					ExerciseName: lift.ExerciseName,
					WeightLifted: lift.WeightLifted,
					Reps:         lift.Reps,
					UserID:       lift.UserID,
					WorkoutID:    lift.WorkoutID,
				}
				store.EXPECT().CreateLift(gomock.Any(), gomock.Eq(args)).Times(0)
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

			url := "/lift"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestCreateLifts(t *testing.T) {
	args, lifts := generateCompleteWorkoutLifts()

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
				"exercise_name": args.Exercisenames,
				"weight":        args.Weights,
				"reps":          args.Reps,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateLifts(gomock.Any(), gomock.Eq(args)).Times(1).Return(lifts, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateLiftsResponse(t, recorder.Body, lifts)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateLifts(gomock.Any(), gomock.Any()).Times(0).Return([]db.Lift{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"exercise_name": args.Exercisenames,
				"weight":        args.Weights,
				"reps":          args.Reps,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateLifts(gomock.Any(), gomock.Eq(args)).Times(1).Return([]db.Lift{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			body: gin.H{
				"exercise_name": args.Exercisenames,
				"weight":        args.Weights,
				"reps":          args.Reps,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateLifts(gomock.Any(), gomock.Eq(args)).Times(0)
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

			url := fmt.Sprintf("/lift/%s/%s", lifts[0].WorkoutID, lifts[0].UserID)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestGetLift(t *testing.T) {
	lift := generateRandLift()

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
				args := db.GetLiftParams{
					UserID: lift.UserID,
					ID:     lift.ID,
				}
				store.EXPECT().GetLift(gomock.Any(), gomock.Eq(args)).Times(1).Return(lift, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateLiftResponse(t, recorder.Body, lift)
			},
		},
		{
			name: "InternalError",
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.GetLiftParams{
					UserID: lift.UserID,
					ID:     lift.ID,
				}
				store.EXPECT().GetLift(gomock.Any(), gomock.Eq(args)).Times(1).Return(db.Lift{}, sql.ErrConnDone)
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
				args := db.GetLiftParams{
					UserID: lift.UserID,
					ID:     lift.ID,
				}
				store.EXPECT().GetLift(gomock.Any(), gomock.Eq(args)).Times(0)
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

			url := fmt.Sprintf("/lift/%s/%s", lift.ID, lift.UserID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestListLifts(t *testing.T) {
	lifts := generateRandLifts()

	type Query struct {
		PageID   int
		PageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				PageSize: 5,
				PageID:   1,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.ListLiftsParams{
					Limit:  5,
					Offset: 0,
					UserID: lifts[0].UserID,
				}
				store.EXPECT().ListLifts(gomock.Any(), gomock.Eq(args)).Times(1).Return(lifts, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateLiftsResponse(t, recorder.Body, lifts)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				PageID:   -1,
				PageSize: 100000,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListLifts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				PageSize: 5,
				PageID:   1,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.ListLiftsParams{
					Limit:  5,
					Offset: 0,
					UserID: lifts[0].UserID,
				}
				store.EXPECT().ListLifts(gomock.Any(), gomock.Eq(args)).Times(1).Return([]db.Lift{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			query: Query{
				PageSize: 5,
				PageID:   1,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.ListLiftsParams{
					Limit:  5,
					Offset: 0,
					UserID: lifts[0].UserID,
				}
				store.EXPECT().ListLifts(gomock.Any(), gomock.Eq(args)).Times(0)
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

			url := fmt.Sprintf("/lift/history/%s", lifts[0].UserID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			qParams := req.URL.Query()
			qParams.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
			qParams.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			req.URL.RawQuery = qParams.Encode()

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestListPrs(t *testing.T) {
	lifts := generateRandLifts()

	type Query struct {
		PageID   int
		PageSize int
		OrderBY  string
		UserID   uuid.UUID
	}

	testCases := []struct {
		name          string
		query         Query
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				PageSize: 5,
				PageID:   1,
				OrderBY:  "weight",
				UserID:   lifts[0].UserID,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPRs(gomock.Any(), gomock.Any()).Times(1).Return(lifts, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateLiftsResponse(t, recorder.Body, lifts)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				PageID:   -1,
				PageSize: 100000,
				OrderBY:  "weight",
				UserID:   lifts[0].UserID,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPRs(gomock.Any(), gomock.Any()).Times(0)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				PageID:   1,
				PageSize: 5,
				OrderBY:  "weight",
				UserID:   lifts[0].UserID,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPRs(gomock.Any(), gomock.Any()).Times(1).Return([]db.Lift{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			query: Query{
				PageID:   1,
				PageSize: 5,
				OrderBY:  "weight",
				UserID:   lifts[0].UserID,
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPRs(gomock.Any(), gomock.Any()).Times(0)
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

			url := fmt.Sprintf("/lift/history/pr/%s/%s", tc.query.OrderBY, tc.query.UserID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			qParams := req.URL.Query()
			qParams.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
			qParams.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			req.URL.RawQuery = qParams.Encode()

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestListPRsByExercise(t *testing.T) {
	lifts := generateRandLifts()

	type Query struct {
		PageID       int
		PageSize     int
		ExerciseName string
		OrderBy      string
		UserID       uuid.UUID
	}

	testCases := []struct {
		name          string
		query         Query
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				PageSize:     5,
				PageID:       1,
				ExerciseName: lifts[0].ExerciseName,
				UserID:       lifts[0].UserID,
				OrderBy:      "weight",
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPRsByExercise(gomock.Any(), gomock.Any()).Times(1).Return(lifts, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				validateLiftsResponse(t, recorder.Body, lifts)
			},
		},
		{
			name: "InternalError",
			query: Query{
				PageSize:     5,
				PageID:       1,
				ExerciseName: lifts[0].ExerciseName,
				UserID:       lifts[0].UserID,
				OrderBy:      "weight",
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPRsByExercise(gomock.Any(), gomock.Any()).Times(1).Return([]db.Lift{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			query: Query{
				PageSize:     5,
				PageID:       1,
				ExerciseName: lifts[0].ExerciseName,
				UserID:       lifts[0].UserID,
				OrderBy:      "weight",
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPRsByExercise(gomock.Any(), gomock.Any()).Times(0)
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

			url := fmt.Sprintf("/lift/pr/%s/%s/%s", tc.query.ExerciseName, tc.query.OrderBy, tc.query.UserID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			qParams := req.URL.Query()
			qParams.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
			qParams.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			req.URL.RawQuery = qParams.Encode()

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func TestListPRsByMuscleGroup(t *testing.T) {
	userId := uuid.New()
	n := 5
	lifts := make([]db.ListPRsByMuscleGroupRow, n)
	for i := 0; i < n; i++ {
		lifts[i] = db.ListPRsByMuscleGroupRow{
			MuscleGroup:  "Chest",
			ExerciseName: util.RandomString(5),
			ID:           uuid.New(),
			WeightLifted: float32(util.RandomInt(100, 200)),
			Reps:         int16(util.RandomInt(5, 12)),
		}
	}

	type Query struct {
		PageID      int
		PageSize    int
		MuscleGroup string
		OrderBy     string
		UserID      uuid.UUID
	}

	testCases := []struct {
		name          string
		query         Query
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				PageSize:    5,
				PageID:      1,
				MuscleGroup: "Chest",
				UserID:      userId,
				OrderBy:     "weight",
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPRsByMuscleGroup(gomock.Any(), gomock.Any()).Times(1).Return(lifts, nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				PageSize:    5,
				PageID:      1,
				MuscleGroup: "Chest",
				UserID:      userId,
				OrderBy:     "weight",
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPRsByMuscleGroup(gomock.Any(), gomock.Any()).Times(1).Return([]db.ListPRsByMuscleGroupRow{}, sql.ErrConnDone)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unauthorized",
			query: Query{
				PageSize:    5,
				PageID:      1,
				MuscleGroup: "Chest",
				UserID:      userId,
				OrderBy:     "weight",
			},
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListPRsByMuscleGroup(gomock.Any(), gomock.Any()).Times(0)
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

			url := fmt.Sprintf("/lift/pr/group/%s/%s/%s", tc.query.MuscleGroup, tc.query.OrderBy, tc.query.UserID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			qParams := req.URL.Query()
			qParams.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
			qParams.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			req.URL.RawQuery = qParams.Encode()

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

// @todo Fix this failing test.
// func TestUpdateLift(t *testing.T) {
// 	lift := generateRandLift()

// 	type Query struct {
// 		PageSize int
// 		PageID   int
// 	}

// 	testCases := []struct {
// 		name          string
// 		query         Query
// 		body          gin.H
// 		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkRes      func(recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"weight_lifted": "20.5",
// 				"reps":          "10",
// 			},
// 			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
// 				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				args := db.UpdateLiftParams{
// 					ID:      lift.ID,
// 					Column1: "20.5",
// 					Column2: "10",
// 				}
// 				store.EXPECT().UpdateLift(gomock.Any(), gomock.Eq(args)).Times(1)
// 			},
// 			checkRes: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			url := fmt.Sprintf("/lift/%s", lift.ID)
// 			req, err := http.NewRequest(http.MethodPatch, url, nil)
// 			require.NoError(t, err)

// 			tc.configureAuth(t, req, server.tokenCreator)
// 			server.router.ServeHTTP(recorder, req)
// 			tc.checkRes(recorder)
// 		})
// 	}
// }

func TestDeleteLift(t *testing.T) {
	lift := generateRandLift()

	testCases := []struct {
		name          string
		liftID        uuid.UUID
		configureAuth func(t *testing.T, request *http.Request, tokenCreator token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkRes      func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK-Deleted",
			liftID: lift.ID,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
				addAuthHeader(t, request, tokenCreator, bearerType, uuid.New(), time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteLift(gomock.Any(), gomock.Eq(lift.ID)).Times(1).Return(nil)
			},
			checkRes: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			liftID: lift.ID,
			configureAuth: func(t *testing.T, request *http.Request, tokenCreator token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteLift(gomock.Any(), gomock.Eq(lift.ID)).Times(0)
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

			url := fmt.Sprintf("/lift/%s", tc.liftID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.configureAuth(t, req, server.tokenCreator)
			server.router.ServeHTTP(recorder, req)
			tc.checkRes(recorder)
		})
	}
}

func generateRandLift() db.Lift {
	return db.Lift{
		ID:           uuid.New(),
		ExerciseName: util.RandomString(5),
		WeightLifted: float32(util.RandomInt(100, 200)),
		Reps:         int16(util.RandomInt(5, 12)),
		UserID:       uuid.New(),
		WorkoutID:    uuid.New(),
	}
}

func generateCompleteWorkoutLifts() (db.CreateLiftsParams, []db.Lift) {

	n := 5
	userID := uuid.New()
	workoutID := uuid.New()
	lifts := make([]db.Lift, n)
	createLiftsArgs := db.CreateLiftsParams{
		Exercisenames: make([]string, n),
		Weights:       make([]float32, n),
		Reps:          make([]int16, n),
		UserID:        make([]uuid.UUID, n),
		WorkoutID:     make([]uuid.UUID, n),
	}

	for i := 0; i < n; i++ {
		createLiftsArgs.Exercisenames[i] = util.RandomString(5)
		createLiftsArgs.Weights[i] = float32(util.RandomInt(100, 220))
		createLiftsArgs.Reps[i] = int16(util.RandomInt(6, 12))
		createLiftsArgs.UserID[i] = userID
		createLiftsArgs.WorkoutID[i] = workoutID

		lifts[i] = db.Lift{
			ExerciseName: createLiftsArgs.Exercisenames[i],
			Reps:         createLiftsArgs.Reps[i],
			WeightLifted: createLiftsArgs.Weights[i],
			WorkoutID:    createLiftsArgs.WorkoutID[i],
			UserID:       createLiftsArgs.UserID[i],
			ID:           uuid.New(),
		}
	}

	return createLiftsArgs, lifts
}

func generateRandLifts() []db.Lift {
	userID := uuid.New()
	workoutID := uuid.New()
	n := 5
	lifts := make([]db.Lift, n)
	for i := 0; i < n; i++ {
		lifts[i] = db.Lift{
			ID:           uuid.New(),
			ExerciseName: util.RandomString(5),
			WeightLifted: float32(util.RandomInt(100, 200)),
			Reps:         int16(util.RandomInt(5, 12)),
			UserID:       userID,
			WorkoutID:    workoutID,
		}
	}
	return lifts
}

func validateLiftResponse(t *testing.T, body *bytes.Buffer, lift db.Lift) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resLift db.Lift
	err = json.Unmarshal(data, &resLift)
	require.NoError(t, err)
	require.Equal(t, lift, resLift)
}

func validateLiftsResponse(t *testing.T, body *bytes.Buffer, lift []db.Lift) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resLifts []db.Lift
	err = json.Unmarshal(data, &resLifts)
	require.NoError(t, err)
	require.Equal(t, lift, resLifts)
}
