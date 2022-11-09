package api

import (
	"net/http"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/gin-gonic/gin"
)

type setUserId struct {
	UserId string `json:"user_id" binding:"required"`
}

type createSetReq struct {
	UserId string `json:"user_id" binding:"required"`
}

func (server *Server) createSet(ctx *gin.Context) {
	var req createSetReq
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
