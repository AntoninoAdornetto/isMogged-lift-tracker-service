package api

import (
	"context"
	"database/sql"
	"net/http"

	db "github.com/AntoninoAdornetto/isMogged-lift-tracker-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createExerciseReq struct {
	Name        string `json:"name" binding:"required,min=3"`
	MuscleGroup string `json:"muscle_group" binding:"required"`
	Category    string `json:"category"`
}

func (server *Server) createExercise(ctx *gin.Context) {
	var req createExerciseReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateExerciseParams{
		Name:        req.Name,
		MuscleGroup: req.MuscleGroup,
		Category:    req.Category,
	}

	ex, err := server.store.CreateExercise(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ex)
}

type getExerciseReq struct {
	Name string `uri:"name" binding:"required"`
}

func (server *Server) getExercise(ctx *gin.Context) {
	var req getExerciseReq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	exercise, err := server.store.GetExercise(ctx, req.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, exercise)
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

type listByMuscleGroupReq struct {
	MuscleGroup string `uri:"muscle_group" binding:"required"`
}

func (server *Server) getMuscleGroupExercises(ctx *gin.Context) {
	var req listByMuscleGroupReq
	var query listExercisesReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListByMuscleGroupParams{
		MuscleGroup: req.MuscleGroup,
		Limit:       query.PageSize,
		Offset:      (query.PageID - 1) * query.PageSize,
	}

	exercises, err := server.store.ListByMuscleGroup(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, exercises)
}

type updateExerciseReq struct {
	Name        string `json:"name"`
	MuscleGroup string `json:"muscle_group"`
	Category    string `json:"category"`
}

func (server *Server) updateExercise(ctx *gin.Context) {
	var req updateExerciseReq
	var uri getExerciseReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateExerciseParams{
		Name:    uri.Name,
		Column1: req.Name,
		Column2: req.MuscleGroup,
		Column3: req.Category,
	}

	patch, err := server.store.UpdateExercise(context.Background(), args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, patch)
}

func (server *Server) deleteExercise(ctx *gin.Context) {
	var req getExerciseReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteExercise(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
