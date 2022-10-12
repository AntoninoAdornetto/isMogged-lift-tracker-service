package api

import (
	"net/http"

	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createExersiseReq struct {
	ExersiseName string `json:"exersise_name" binding:"required,min=3"`
	MuscleGroup  string `json:"muscle_group" binding:"required,min=3"`
}

func (server *Server) createExersise(ctx *gin.Context) {
	var req createExersiseReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateExersiseParams{
		ExersiseName: req.ExersiseName,
		MuscleGroup:  req.MuscleGroup,
	}

	ex, err := server.store.CreateExersise(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ex)
}

type getExersiseReq struct {
	ExersiseName string `uri:"exersise_name" binding:"required"`
}

func (server *Server) getExersise(ctx *gin.Context) {
	var req getExersiseReq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ex, err := server.store.GetExersise(ctx, req.ExersiseName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ex)
}

type listExersisesReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
}

func (server *Server) listExersises(ctx *gin.Context) {
	var req listExersisesReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListExersisesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	exs, err := server.store.ListExersises(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, exs)
}

type getMuscleGroupExersisesReq struct {
	MuscleGroup string `uri:"muscle_group" binding:"required"`
}

func (server *Server) getMuscleGroupExersises(ctx *gin.Context) {
	var req getMuscleGroupExersisesReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	exs, err := server.store.GetMuscleGroupExersises(ctx, req.MuscleGroup)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, exs)
}

type updateExersiseNameReq struct {
	CurrentName string `json:"current_name" binding:"required"`
	NewName     string `json:"new_name" binding:"required"`
}

func (server *Server) updateExersiseName(ctx *gin.Context) {
	var req updateExersiseNameReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateExersiseNameParams{
		ExersiseName:   req.NewName,
		ExersiseName_2: req.CurrentName,
	}

	err := server.store.UpdateExersiseName(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

type updateExersiseMuscleGroupReq struct {
	MuscleGroup  string `json:"muscle_group" binding:"required"`
	ExersiseName string `json:"exersise_name" binding:"required"`
}

func (server *Server) updateExersiseMuscleGroup(ctx *gin.Context) {
	var req updateExersiseMuscleGroupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateExersiseMuscleGroupParams{
		MuscleGroup:  req.MuscleGroup,
		ExersiseName: req.ExersiseName,
	}

	err := server.store.UpdateExersiseMuscleGroup(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
