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
