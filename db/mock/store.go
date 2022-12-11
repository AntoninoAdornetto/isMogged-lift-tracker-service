// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/AntoninoAdornetto/lift_tracker/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockStore) CreateAccount(arg0 context.Context, arg1 db.CreateAccountParams) (db.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0, arg1)
	ret0, _ := ret[0].(db.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockStoreMockRecorder) CreateAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockStore)(nil).CreateAccount), arg0, arg1)
}

// CreateCategory mocks base method.
func (m *MockStore) CreateCategory(arg0 context.Context, arg1 string) (db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCategory", arg0, arg1)
	ret0, _ := ret[0].(db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCategory indicates an expected call of CreateCategory.
func (mr *MockStoreMockRecorder) CreateCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCategory", reflect.TypeOf((*MockStore)(nil).CreateCategory), arg0, arg1)
}

// CreateExercise mocks base method.
func (m *MockStore) CreateExercise(arg0 context.Context, arg1 db.CreateExerciseParams) (db.Exercise, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateExercise", arg0, arg1)
	ret0, _ := ret[0].(db.Exercise)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateExercise indicates an expected call of CreateExercise.
func (mr *MockStoreMockRecorder) CreateExercise(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateExercise", reflect.TypeOf((*MockStore)(nil).CreateExercise), arg0, arg1)
}

// CreateLift mocks base method.
func (m *MockStore) CreateLift(arg0 context.Context, arg1 db.CreateLiftParams) (db.Lift, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLift", arg0, arg1)
	ret0, _ := ret[0].(db.Lift)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLift indicates an expected call of CreateLift.
func (mr *MockStoreMockRecorder) CreateLift(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLift", reflect.TypeOf((*MockStore)(nil).CreateLift), arg0, arg1)
}

// CreateMuscleGroup mocks base method.
func (m *MockStore) CreateMuscleGroup(arg0 context.Context, arg1 string) (db.MuscleGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMuscleGroup", arg0, arg1)
	ret0, _ := ret[0].(db.MuscleGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMuscleGroup indicates an expected call of CreateMuscleGroup.
func (mr *MockStoreMockRecorder) CreateMuscleGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMuscleGroup", reflect.TypeOf((*MockStore)(nil).CreateMuscleGroup), arg0, arg1)
}

// CreateWorkout mocks base method.
func (m *MockStore) CreateWorkout(arg0 context.Context, arg1 db.CreateWorkoutParams) (db.Workout, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWorkout", arg0, arg1)
	ret0, _ := ret[0].(db.Workout)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWorkout indicates an expected call of CreateWorkout.
func (mr *MockStoreMockRecorder) CreateWorkout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWorkout", reflect.TypeOf((*MockStore)(nil).CreateWorkout), arg0, arg1)
}

// DeleteAccount mocks base method.
func (m *MockStore) DeleteAccount(arg0 context.Context, arg1 uuid.UUID) (db.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", arg0, arg1)
	ret0, _ := ret[0].(db.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockStoreMockRecorder) DeleteAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockStore)(nil).DeleteAccount), arg0, arg1)
}

// DeleteCategory mocks base method.
func (m *MockStore) DeleteCategory(arg0 context.Context, arg1 int16) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCategory", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCategory indicates an expected call of DeleteCategory.
func (mr *MockStoreMockRecorder) DeleteCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCategory", reflect.TypeOf((*MockStore)(nil).DeleteCategory), arg0, arg1)
}

// DeleteExercise mocks base method.
func (m *MockStore) DeleteExercise(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteExercise", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteExercise indicates an expected call of DeleteExercise.
func (mr *MockStoreMockRecorder) DeleteExercise(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExercise", reflect.TypeOf((*MockStore)(nil).DeleteExercise), arg0, arg1)
}

// DeleteGroup mocks base method.
func (m *MockStore) DeleteGroup(arg0 context.Context, arg1 string) (db.MuscleGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroup", arg0, arg1)
	ret0, _ := ret[0].(db.MuscleGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteGroup indicates an expected call of DeleteGroup.
func (mr *MockStoreMockRecorder) DeleteGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockStore)(nil).DeleteGroup), arg0, arg1)
}

// DeleteLift mocks base method.
func (m *MockStore) DeleteLift(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLift", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLift indicates an expected call of DeleteLift.
func (mr *MockStoreMockRecorder) DeleteLift(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLift", reflect.TypeOf((*MockStore)(nil).DeleteLift), arg0, arg1)
}

// DeleteWorkout mocks base method.
func (m *MockStore) DeleteWorkout(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWorkout", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWorkout indicates an expected call of DeleteWorkout.
func (mr *MockStoreMockRecorder) DeleteWorkout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWorkout", reflect.TypeOf((*MockStore)(nil).DeleteWorkout), arg0, arg1)
}

// GetAccount mocks base method.
func (m *MockStore) GetAccount(arg0 context.Context, arg1 uuid.UUID) (db.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0, arg1)
	ret0, _ := ret[0].(db.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockStoreMockRecorder) GetAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockStore)(nil).GetAccount), arg0, arg1)
}

// GetCategory mocks base method.
func (m *MockStore) GetCategory(arg0 context.Context, arg1 int16) (db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategory", arg0, arg1)
	ret0, _ := ret[0].(db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategory indicates an expected call of GetCategory.
func (mr *MockStoreMockRecorder) GetCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategory", reflect.TypeOf((*MockStore)(nil).GetCategory), arg0, arg1)
}

// GetExercise mocks base method.
func (m *MockStore) GetExercise(arg0 context.Context, arg1 string) (db.Exercise, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExercise", arg0, arg1)
	ret0, _ := ret[0].(db.Exercise)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetExercise indicates an expected call of GetExercise.
func (mr *MockStoreMockRecorder) GetExercise(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExercise", reflect.TypeOf((*MockStore)(nil).GetExercise), arg0, arg1)
}

// GetLift mocks base method.
func (m *MockStore) GetLift(arg0 context.Context, arg1 db.GetLiftParams) (db.Lift, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLift", arg0, arg1)
	ret0, _ := ret[0].(db.Lift)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLift indicates an expected call of GetLift.
func (mr *MockStoreMockRecorder) GetLift(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLift", reflect.TypeOf((*MockStore)(nil).GetLift), arg0, arg1)
}

// GetMuscleGroup mocks base method.
func (m *MockStore) GetMuscleGroup(arg0 context.Context, arg1 string) (db.MuscleGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMuscleGroup", arg0, arg1)
	ret0, _ := ret[0].(db.MuscleGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMuscleGroup indicates an expected call of GetMuscleGroup.
func (mr *MockStoreMockRecorder) GetMuscleGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMuscleGroup", reflect.TypeOf((*MockStore)(nil).GetMuscleGroup), arg0, arg1)
}

// GetMuscleGroups mocks base method.
func (m *MockStore) GetMuscleGroups(arg0 context.Context) ([]db.MuscleGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMuscleGroups", arg0)
	ret0, _ := ret[0].([]db.MuscleGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMuscleGroups indicates an expected call of GetMuscleGroups.
func (mr *MockStoreMockRecorder) GetMuscleGroups(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMuscleGroups", reflect.TypeOf((*MockStore)(nil).GetMuscleGroups), arg0)
}

// GetWorkout mocks base method.
func (m *MockStore) GetWorkout(arg0 context.Context, arg1 uuid.UUID) ([]db.GetWorkoutRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWorkout", arg0, arg1)
	ret0, _ := ret[0].([]db.GetWorkoutRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWorkout indicates an expected call of GetWorkout.
func (mr *MockStoreMockRecorder) GetWorkout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWorkout", reflect.TypeOf((*MockStore)(nil).GetWorkout), arg0, arg1)
}

// ListAccounts mocks base method.
func (m *MockStore) ListAccounts(arg0 context.Context, arg1 db.ListAccountsParams) ([]db.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAccounts", arg0, arg1)
	ret0, _ := ret[0].([]db.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAccounts indicates an expected call of ListAccounts.
func (mr *MockStoreMockRecorder) ListAccounts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAccounts", reflect.TypeOf((*MockStore)(nil).ListAccounts), arg0, arg1)
}

// ListByMuscleGroup mocks base method.
func (m *MockStore) ListByMuscleGroup(arg0 context.Context, arg1 db.ListByMuscleGroupParams) ([]db.Exercise, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByMuscleGroup", arg0, arg1)
	ret0, _ := ret[0].([]db.Exercise)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByMuscleGroup indicates an expected call of ListByMuscleGroup.
func (mr *MockStoreMockRecorder) ListByMuscleGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByMuscleGroup", reflect.TypeOf((*MockStore)(nil).ListByMuscleGroup), arg0, arg1)
}

// ListCategories mocks base method.
func (m *MockStore) ListCategories(arg0 context.Context) ([]db.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCategories", arg0)
	ret0, _ := ret[0].([]db.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCategories indicates an expected call of ListCategories.
func (mr *MockStoreMockRecorder) ListCategories(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCategories", reflect.TypeOf((*MockStore)(nil).ListCategories), arg0)
}

// ListExercises mocks base method.
func (m *MockStore) ListExercises(arg0 context.Context, arg1 db.ListExercisesParams) ([]db.Exercise, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListExercises", arg0, arg1)
	ret0, _ := ret[0].([]db.Exercise)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListExercises indicates an expected call of ListExercises.
func (mr *MockStoreMockRecorder) ListExercises(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListExercises", reflect.TypeOf((*MockStore)(nil).ListExercises), arg0, arg1)
}

// ListLifts mocks base method.
func (m *MockStore) ListLifts(arg0 context.Context, arg1 db.ListLiftsParams) ([]db.Lift, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListLifts", arg0, arg1)
	ret0, _ := ret[0].([]db.Lift)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListLifts indicates an expected call of ListLifts.
func (mr *MockStoreMockRecorder) ListLifts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListLifts", reflect.TypeOf((*MockStore)(nil).ListLifts), arg0, arg1)
}

// ListPRs mocks base method.
func (m *MockStore) ListPRs(arg0 context.Context, arg1 db.ListPRsParams) ([]db.Lift, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPRs", arg0, arg1)
	ret0, _ := ret[0].([]db.Lift)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPRs indicates an expected call of ListPRs.
func (mr *MockStoreMockRecorder) ListPRs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPRs", reflect.TypeOf((*MockStore)(nil).ListPRs), arg0, arg1)
}

// ListPRsByExercise mocks base method.
func (m *MockStore) ListPRsByExercise(arg0 context.Context, arg1 db.ListPRsByExerciseParams) ([]db.Lift, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPRsByExercise", arg0, arg1)
	ret0, _ := ret[0].([]db.Lift)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPRsByExercise indicates an expected call of ListPRsByExercise.
func (mr *MockStoreMockRecorder) ListPRsByExercise(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPRsByExercise", reflect.TypeOf((*MockStore)(nil).ListPRsByExercise), arg0, arg1)
}

// ListPRsByMuscleGroup mocks base method.
func (m *MockStore) ListPRsByMuscleGroup(arg0 context.Context, arg1 db.ListPRsByMuscleGroupParams) ([]db.ListPRsByMuscleGroupRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPRsByMuscleGroup", arg0, arg1)
	ret0, _ := ret[0].([]db.ListPRsByMuscleGroupRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPRsByMuscleGroup indicates an expected call of ListPRsByMuscleGroup.
func (mr *MockStoreMockRecorder) ListPRsByMuscleGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPRsByMuscleGroup", reflect.TypeOf((*MockStore)(nil).ListPRsByMuscleGroup), arg0, arg1)
}

// ListWorkouts mocks base method.
func (m *MockStore) ListWorkouts(arg0 context.Context, arg1 db.ListWorkoutsParams) ([]db.Workout, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWorkouts", arg0, arg1)
	ret0, _ := ret[0].([]db.Workout)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWorkouts indicates an expected call of ListWorkouts.
func (mr *MockStoreMockRecorder) ListWorkouts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWorkouts", reflect.TypeOf((*MockStore)(nil).ListWorkouts), arg0, arg1)
}

// UpdateCategory mocks base method.
func (m *MockStore) UpdateCategory(arg0 context.Context, arg1 db.UpdateCategoryParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCategory", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCategory indicates an expected call of UpdateCategory.
func (mr *MockStoreMockRecorder) UpdateCategory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCategory", reflect.TypeOf((*MockStore)(nil).UpdateCategory), arg0, arg1)
}

// UpdateExerciseName mocks base method.
func (m *MockStore) UpdateExerciseName(arg0 context.Context, arg1 db.UpdateExerciseNameParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateExerciseName", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateExerciseName indicates an expected call of UpdateExerciseName.
func (mr *MockStoreMockRecorder) UpdateExerciseName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateExerciseName", reflect.TypeOf((*MockStore)(nil).UpdateExerciseName), arg0, arg1)
}

// UpdateFinishTime mocks base method.
func (m *MockStore) UpdateFinishTime(arg0 context.Context, arg1 db.UpdateFinishTimeParams) (db.Workout, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFinishTime", arg0, arg1)
	ret0, _ := ret[0].(db.Workout)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFinishTime indicates an expected call of UpdateFinishTime.
func (mr *MockStoreMockRecorder) UpdateFinishTime(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFinishTime", reflect.TypeOf((*MockStore)(nil).UpdateFinishTime), arg0, arg1)
}

// UpdateGroup mocks base method.
func (m *MockStore) UpdateGroup(arg0 context.Context, arg1 db.UpdateGroupParams) (db.MuscleGroup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroup", arg0, arg1)
	ret0, _ := ret[0].(db.MuscleGroup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateGroup indicates an expected call of UpdateGroup.
func (mr *MockStoreMockRecorder) UpdateGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroup", reflect.TypeOf((*MockStore)(nil).UpdateGroup), arg0, arg1)
}

// UpdateLift mocks base method.
func (m *MockStore) UpdateLift(arg0 context.Context, arg1 db.UpdateLiftParams) (db.Lift, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLift", arg0, arg1)
	ret0, _ := ret[0].(db.Lift)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateLift indicates an expected call of UpdateLift.
func (mr *MockStoreMockRecorder) UpdateLift(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLift", reflect.TypeOf((*MockStore)(nil).UpdateLift), arg0, arg1)
}

// UpdateMuscleGroup mocks base method.
func (m *MockStore) UpdateMuscleGroup(arg0 context.Context, arg1 db.UpdateMuscleGroupParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMuscleGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMuscleGroup indicates an expected call of UpdateMuscleGroup.
func (mr *MockStoreMockRecorder) UpdateMuscleGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMuscleGroup", reflect.TypeOf((*MockStore)(nil).UpdateMuscleGroup), arg0, arg1)
}

// UpdateWeight mocks base method.
func (m *MockStore) UpdateWeight(arg0 context.Context, arg1 db.UpdateWeightParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateWeight", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateWeight indicates an expected call of UpdateWeight.
func (mr *MockStoreMockRecorder) UpdateWeight(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateWeight", reflect.TypeOf((*MockStore)(nil).UpdateWeight), arg0, arg1)
}
