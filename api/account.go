package api

import (
	"net/http"
	"time"

	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountReq struct {
	Lifter string `json:"lifter" binding:"required"`
	// BirthDate string `json:"birth_date" binding:"required"`
	Weight int32 `json:"weight" binding:"required"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: Fix JSON encoded date parsing.
	// bDay, err := time.Parse("1999-09-05", req.BirthDate)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

	arg := db.CreateAccountParams{
		Lifter:    req.Lifter,
		BirthDate: time.Now(),
		Weight:    req.Weight,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
