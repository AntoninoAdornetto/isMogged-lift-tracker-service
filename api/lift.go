package api

import (
	"net/http"

	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/gin-gonic/gin"
)

type liftUserId struct {
	UserId string `json:"user_id" binding:"required"`
}

type liftPaginationReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
}

type createLiftReq struct {
	ExersiseName string  `json:"exercise_name" binding:"required"`
	Weight       float32 `json:"weight" binding:"required"`
	Reps         int32   `json:"reps" binding:"required"`
	UserId       string  `json:"user_id" binding:"required"`
	SetId        string  `json:"set_id" binding:"required"`
}

func (server *Server) createLift(ctx *gin.Context) {
	var req createLiftReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userId, err := util.ParseUUIDStr(req.UserId, ctx)
	if err != nil {
		return
	}

	setId, err := util.ParseUUIDStr(req.SetId, ctx)
	if err != nil {
		return
	}

	args := db.CreateLiftParams{
		ExerciseName: req.ExersiseName,
		Weight:       req.Weight,
		Reps:         req.Reps,
		UserID:       userId,
		SetID:        setId,
	}

	lift, err := server.store.CreateLift(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, lift)
}

type getLiftReq struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getLift(ctx *gin.Context) {
	var req getLiftReq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	lift, err := server.store.GetLift(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, lift)
}

func (server *Server) listLifts(ctx *gin.Context) {
	var acc liftUserId
	var req liftPaginationReq

	if err := ctx.BindJSON(&acc); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := util.ParseUUIDStr(acc.UserId, ctx)
	if err != nil {
		return
	}

	args := db.ListLiftsParams{
		UserID: id,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	lifts, err := server.store.ListLifts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, error(err))
		return
	}

	ctx.JSON(http.StatusOK, lifts)
}

func (server *Server) listWeightPRs(ctx *gin.Context) {
	var acc liftUserId
	var req liftPaginationReq

	if err := ctx.BindJSON(&acc); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := util.ParseUUIDStr(acc.UserId, ctx)
	if err != nil {
		return
	}

	args := db.ListWeightPRLiftsParams{
		UserID: id,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	lifts, err := server.store.ListWeightPRLifts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, error(err))
		return
	}

	ctx.JSON(http.StatusOK, lifts)
}

// -- name: ListNamedLiftWeightPRs :many
// SELECT * FROM lift
// WHERE user_id = $1 AND exercise_name = $2
// ORDER BY weight DESC
// LIMIT $2
// OFFSET $3;

type listNamedLiftWeightPRReq struct {
	ExerciseName string `uri:"exercise_name" binding:"required"`
}

func (server *Server) listNamedLiftWeightPRs(ctx *gin.Context) {
	var acc liftUserId
	var req listNamedLiftWeightPRReq
	var pagination liftPaginationReq

	if err := ctx.BindJSON(&acc); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.BindQuery(&pagination); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := util.ParseUUIDStr(acc.UserId, ctx)
	if err != nil {
		return
	}

	args := db.ListNamedLiftWeightPRsParams{
		UserID:       id,
		ExerciseName: req.ExerciseName,
		Limit:        pagination.PageSize,
		Offset:       (pagination.PageID - 1) * pagination.PageSize,
	}

	lifts, err := server.store.ListNamedLiftWeightPRs(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, error(err))
		return
	}

	ctx.JSON(http.StatusOK, lifts)
}
