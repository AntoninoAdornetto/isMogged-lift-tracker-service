package api

import (
	"database/sql"
	"net/http"

	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createExerciseReq struct {
	ExerciseName string `json:"exercise_name" binding:"required,min=3"`
	MuscleGroup  string `json:"muscle_group" binding:"required,min=3"`
}

func (server *Server) createExercise(ctx *gin.Context) {
	var req createExerciseReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateExerciseParams{
		ExerciseName: req.ExerciseName,
		MuscleGroup:  req.MuscleGroup,
	}

	ex, err := server.store.CreateExercise(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ex)
}

type getExerciseReq struct {
	ExerciseName string `uri:"exercise_name" binding:"required"`
}

func (server *Server) getExercise(ctx *gin.Context) {
	var req getExerciseReq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ex, err := server.store.GetExercise(ctx, req.ExerciseName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ex)
}

type listExercisesReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
}

func (server *Server) listExercises(ctx *gin.Context) {
	var req listExercisesReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListExercisesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	exs, err := server.store.ListExercises(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, exs)
}

type getMuscleGroupExercisesReq struct {
	MuscleGroup string `uri:"muscle_group" binding:"required"`
}

func (server *Server) getMuscleGroupExercises(ctx *gin.Context) {
	var req getMuscleGroupExercisesReq
	var query listExercisesReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.GetMuscleGroupExercisesParams{
		MuscleGroup: req.MuscleGroup,
		Limit:       query.PageSize,
		Offset:      (query.PageID - 1) * query.PageSize,
	}

	exercises, err := server.store.GetMuscleGroupExercises(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, exercises)
}

type updateExerciseNameReq struct {
	CurrentName string `json:"current_name" binding:"required"`
	NewName     string `json:"new_name" binding:"required"`
}

func (server *Server) updateExerciseName(ctx *gin.Context) {
	var req updateExerciseNameReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateExerciseNameParams{
		ExerciseName:   req.NewName,
		ExerciseName_2: req.CurrentName,
	}

	err := server.store.UpdateExerciseName(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

type updateExerciseMuscleGroupReq struct {
	MuscleGroup  string `json:"muscle_group" binding:"required"`
	ExerciseName string `json:"exercise_name" binding:"required"`
}

func (server *Server) updateExerciseMuscleGroup(ctx *gin.Context) {
	var req updateExerciseMuscleGroupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateExerciseMuscleGroupParams{
		MuscleGroup:  req.MuscleGroup,
		ExerciseName: req.ExerciseName,
	}

	err := server.store.UpdateExerciseMuscleGroup(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

type deleteExerciseReq struct {
	ExerciseName string `uri:"exercise_name" binding:"required"`
}

func (server *Server) deleteExercise(ctx *gin.Context) {
	var req deleteExerciseReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteExercise(ctx, req.ExerciseName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
