package api

import (
	"context"
	"database/sql"
	"net/http"

	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createWorkoutReq struct {
	UserId string `uri:"user_id" binding:"required"`
}

func (server *Server) createWorkout(ctx *gin.Context) {
	var req createWorkoutReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user_id, err := uuid.Parse(req.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	workout, err := server.store.CreateWorkout(ctx, user_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, workout)
}

type getWorkoutReq struct {
	UserId    string `uri:"user_id" binding:"required"`
	WorkoutId string `uri:"workout_id" binding:"required"`
}

func (server *Server) getWorkout(ctx *gin.Context) {
	var req getWorkoutReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	workoutId, err := uuid.Parse(req.WorkoutId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.GetWorkoutParams{
		UserID: userId,
		ID:     workoutId,
	}

	lifts, err := server.store.GetWorkout(context.Background(), args)
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

// type updateWorkoutReq struct {
// 	WorkoutId string `json:"workout_id" binding:"required"`
// }

// type updateDurationReq struct {
// 	finish_time string `json:"finish_time" binding:"required"`
// }

// func (server *Server) updateDurationEnd(ctx *gin.Context) {
// 	var uri updateWorkoutReq
// 	if err := ctx.ShouldBindUri(&uri); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	workoutId, err := uuid.Parse(uri.WorkoutId)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	workout, err := server.store.UpdateDurationEnd(context.Background())
// }

// @todo
// updateDurationEnd
// ListWorkouts
// DeleteWorkout
