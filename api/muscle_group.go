package api

import (
	"database/sql"
	"net/http"
	"strings"

	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createMuscleGroupReq struct {
	GroupName *string `json:"group_name" binding:"required,min=3"`
}

func (server *Server) createMuscleGroup(ctx *gin.Context) {
	var req createMuscleGroupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	mg, err := server.store.CreateMuscleGroup(ctx, *req.GroupName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mg)
}

type getMuscleGroupReq struct {
	GroupName *string `uri:"group_name" binding:"required"`
}

func (server *Server) getMuscleGroup(ctx *gin.Context) {
	var req getMuscleGroupReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	mg, err := server.store.GetMuscleGroup(ctx, *req.GroupName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mg)
}

func (server *Server) listMuscleGroups(ctx *gin.Context) {
	mgs, err := server.store.GetMuscleGroups(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, mgs)
}

type updateMuscleGroupReq struct {
	UpdatedName *string `json:"updated_name" binding:"required"`
	GroupName   *string `json:"group_name" binding:"required"`
}

func (server *Server) updateMuscleGroup(ctx *gin.Context) {
	var req updateMuscleGroupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateGroupParams{
		GroupName:   strings.ToLower(*req.UpdatedName),
		GroupName_2: *req.GroupName,
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

type deleteMuscleGroupReq struct {
	GroupName *string `uri:"group_name" binding:"required"`
}

func (server *Server) deleteMuscleGroup(ctx *gin.Context) {
	var req deleteMuscleGroupReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err := server.store.DeleteGroup(ctx, *req.GroupName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, err)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
