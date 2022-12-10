package api

import (
	"context"
	"net/http"
	"strconv"

	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type liftPaginationReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=50"`
}

type createLiftReq struct {
	ExersiseName string  `json:"exercise_name" binding:"required"`
	Weight       float32 `json:"weight" binding:"required"`
	Reps         int32   `json:"reps" binding:"required"`
	UserId       string  `json:"user_id" binding:"required"`
	WorkoutID    string  `json:"workout_id" binding:"required"`
}

func (server *Server) createLift(ctx *gin.Context) {
	var req createLiftReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	workoutId, err := uuid.Parse(req.WorkoutID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateLiftParams{
		ExerciseName: req.ExersiseName,
		WeightLifted: req.Weight,
		Reps:         int16(req.Reps),
		UserID:       userId,
		WorkoutID:    workoutId,
	}

	lift, err := server.store.CreateLift(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, lift)
}

type getLiftReq struct {
	ID     string `uri:"id" binding:"required"`
	UserID string `uri:"user_id" binding:"required"`
}

func (server *Server) getLift(ctx *gin.Context) {
	var req getLiftReq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	liftId, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	lift, err := server.store.GetLift(ctx, db.GetLiftParams{
		UserID: userId,
		ID:     liftId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, lift)
}

type listLiftsReq struct {
	UserID string `uri:"user_id" binding:"required"`
}

func (server *Server) listLifts(ctx *gin.Context) {
	var uri listLiftsReq
	var req liftPaginationReq

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(uri.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListLiftsParams{
		UserID: userId,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	lifts, err := server.store.ListLifts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, error(err))
		return
	}

	ctx.JSON(http.StatusOK, lifts)
}

type listPRsReq struct {
	UserID  string `uri:"user_id" binding:"required"`
	OrderBy string `uri:"order_by" binding:"required"`
}

func (server *Server) listPRs(ctx *gin.Context) {
	var uri listPRsReq
	var req liftPaginationReq

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(uri.UserID)
	if err != nil {
		return
	}

	args := db.ListPRsParams{
		UserID:  id,
		Column2: uri.OrderBy,
		Column3: uri.OrderBy,
		Limit:   req.PageSize,
		Offset:  (req.PageID - 1) * req.PageSize,
	}

	lifts, err := server.store.ListPRs(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, error(err))
		return
	}

	ctx.JSON(http.StatusOK, lifts)
}

type listPRsByExerciseReq struct {
	ExerciseName string `uri:"exercise_name" binding:"required"`
	OrderBy      string `uri:"order_by" binding:"required"`
	UserID       string `uri:"user_id" binding:"required"`
}

func (server *Server) listPRsByExercise(ctx *gin.Context) {
	var req listPRsByExerciseReq
	var query liftPaginationReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	lifts, err := server.store.ListPRsByExercise(context.Background(), db.ListPRsByExerciseParams{
		UserID:       userId,
		ExerciseName: req.ExerciseName,
		Limit:        query.PageSize,
		Offset:       (query.PageID - 1) * query.PageSize,
		Column3:      req.OrderBy,
		Column4:      req.OrderBy,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, lifts)
}

type listPRsByMuscleGroupReq struct {
	UserID      string `uri:"user_id" binding:"required"`
	MuscleGroup string `uri:"muscle_group" binding:"required"`
	OrderBy     string `uri:"order_by" binding:"required"`
}

func (server *Server) listPRsByMuscleGroup(ctx *gin.Context) {
	var req listPRsByMuscleGroupReq
	var query liftPaginationReq

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	lifts, err := server.store.ListPRsByMuscleGroup(context.Background(), db.ListPRsByMuscleGroupParams{
		MuscleGroup: req.MuscleGroup,
		UserID:      userId,
		Column3:     req.OrderBy,
		Column4:     req.OrderBy,
		Limit:       query.PageSize,
		Offset:      (query.PageID - 1) * query.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, lifts)
}

type updateLiftReq struct {
	WeightLifted string `json:"weight_lifted"`
	Reps         string `json:"reps"`
}

type getLiftByIdReq struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) updateLift(ctx *gin.Context) {
	var uri getLiftByIdReq
	var req updateLiftReq
	var args db.UpdateLiftParams
	weightPrecision := 32

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	args.ID = id

	patchedWeight, err := strconv.ParseFloat(req.WeightLifted, weightPrecision)
	if err == nil {
		args.Column1 = float32(patchedWeight)
	}

	patchedReps, err := strconv.Atoi(req.Reps)
	if err == nil {
		args.Column2 = int16(patchedReps)
	}

	patched, err := server.store.UpdateLift(context.Background(), args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, patched)
}

func (server *Server) deleteLift(ctx *gin.Context) {
	var lift getLiftByIdReq

	if err := ctx.BindUri(&lift); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(lift.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteLift(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
