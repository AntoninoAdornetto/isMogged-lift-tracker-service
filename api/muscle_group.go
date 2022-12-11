package api

import (
	"database/sql"
	"net/http"
	"strings"

	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createMuscleGroupReq struct {
	Name string `json:"name" binding:"required,min=3"`
}

func (server *Server) createMuscleGroup(ctx *gin.Context) {
	var req createMuscleGroupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	mg, err := server.store.CreateMuscleGroup(ctx, strings.ToLower(req.Name))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mg)
}

type getMuscleGroupReq struct {
	Name string `uri:"name" binding:"required"`
}

func (server *Server) getMuscleGroup(ctx *gin.Context) {
	var req getMuscleGroupReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	muscleGroup, err := server.store.GetMuscleGroup(ctx, req.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, muscleGroup)
}

func (server *Server) listMuscleGroups(ctx *gin.Context) {
	muscleGroups, err := server.store.GetMuscleGroups(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, muscleGroups)
}

type updateMuscleGroupReq struct {
	GroupName   string `json:"name" binding:"required"`
	UpdatedName string `json:"updated_name" binding:"required"`
}

func (server *Server) updateMuscleGroup(ctx *gin.Context) {
	var req updateMuscleGroupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateGroupParams{
		Name:   strings.ToLower(req.UpdatedName),
		Name_2: strings.ToLower(req.GroupName),
	}

	patch, err := server.store.UpdateGroup(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, patch)
}

func (server *Server) deleteMuscleGroup(ctx *gin.Context) {
	var req getMuscleGroupReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	d, err := server.store.DeleteGroup(ctx, req.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, d)
}
