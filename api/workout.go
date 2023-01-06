package api

import (
	"context"
	"database/sql"
	"net/http"

	db "github.com/AntoninoAdornetto/isMogged-lift-tracker-service/db/sqlc"
	"github.com/AntoninoAdornetto/isMogged-lift-tracker-service/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createWorkoutReq struct {
	StartTime int64 `json:"start_time" binding:"required"`
}

type getUserIdReq struct {
	UserId string `uri:"user_id" binding:"required"`
}

func (server *Server) createWorkout(ctx *gin.Context) {
	var uri getUserIdReq
	var req createWorkoutReq
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(uri.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	startTime := util.FormatMSEpoch(req.StartTime)

	workout, err := server.store.CreateWorkout(ctx, db.CreateWorkoutParams{
		UserID:    userId,
		StartTime: startTime,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, workout)
}

type getWorkoutReq struct {
	WorkoutId string `uri:"workout_id" binding:"required"`
}

func (server *Server) getWorkout(ctx *gin.Context) {
	var req getWorkoutReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	workoutId, err := uuid.Parse(req.WorkoutId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	lifts, err := server.store.GetWorkout(context.Background(), workoutId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, lifts)
}

type updateWorkoutReq struct {
	WorkoutId string `uri:"workout_id" binding:"required"`
}

type updateDurationReq struct {
	FinishTime int64 `json:"finish_time" binding:"required"`
}

func (server *Server) updateFinishTime(ctx *gin.Context) {
	var uri updateWorkoutReq
	var req updateDurationReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	workoutId, err := uuid.Parse(uri.WorkoutId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	endTime := util.FormatMSEpoch(req.FinishTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateFinishTimeParams{
		ID:         workoutId,
		FinishTime: endTime,
	}

	workout, err := server.store.UpdateFinishTime(context.Background(), args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, workout)
}

type getWorkoutPagination struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
}

func (server *Server) listWorkouts(ctx *gin.Context) {
	var uri getUserIdReq
	var req getWorkoutPagination
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(uri.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListWorkoutsParams{
		UserID: userId,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	workouts, err := server.store.ListWorkouts(context.Background(), args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, workouts)
}

func (server *Server) deleteWorkout(ctx *gin.Context) {
	var req getWorkoutReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.WorkoutId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteWorkout(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}