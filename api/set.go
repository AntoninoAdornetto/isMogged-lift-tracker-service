package api

import (
	"database/sql"
	"net/http"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/gin-gonic/gin"
)

type getSetReq struct {
	SetId string `uri:"set_id" binding:"required"`
}

type setUserId struct {
	UserId string `json:"user_id" binding:"required"`
}

func (server *Server) createSet(ctx *gin.Context) {
	var req setUserId
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := util.ParseUUIDStr(req.UserId, ctx)
	if err != nil {
		return
	}

	set, err := server.store.CreateSet(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, set)
}

func (server *Server) getSet(ctx *gin.Context) {
	var req getSetReq

	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	setId, err := util.ParseUUIDStr(req.SetId, ctx)
	if err != nil {
		return
	}

	set, err := server.store.GetSet(ctx, setId)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, set)
}

func (server *Server) getLiftSet(ctx *gin.Context) {
	var req getSetReq

	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	setId, err := util.ParseUUIDStr(req.SetId, ctx)
	if err != nil {
		return
	}

	set, err := server.store.GetLiftSets(ctx, setId)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, set)
}

func (server *Server) deleteSet(ctx *gin.Context) {
	var req getSetReq

	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	setId, err := util.ParseUUIDStr(req.SetId, ctx)
	if err != nil {
		return
	}

	_, err = server.store.GetSet(ctx, setId)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.store.DeleteSet(ctx, setId)

	ctx.JSON(http.StatusNoContent, nil)
}
