package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createAccountReq struct {
	Name     string  `json:"name" binding:"required,min=3"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Weight   float32 `json:"weight"`
	BodyFat  float32 `json:"body_fat"`
}

type accountResp struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Weight    float32   `json:"weight"`
	BodyFat   float32   `json:"body_fat"`
	StartDate time.Time `json:"start_date"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateAccountParams{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		Weight:    req.Weight,
		BodyFat:   req.BodyFat,
		StartDate: util.FormatMSEpoch(time.Now().UnixMilli()),
	}

	account, err := server.store.CreateAccount(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
			log.Println(pqErr.Code.Name())
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := accountResp{
		ID:        account.ID,
		Name:      account.Name,
		StartDate: account.StartDate,
		BodyFat:   account.BodyFat,
		Weight:    account.Weight,
		Email:     account.Email,
	}

	ctx.JSON(http.StatusOK, res)
}

type getAccountReq struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := accountResp{
		ID:        account.ID,
		Name:      account.Name,
		StartDate: account.StartDate,
		BodyFat:   account.BodyFat,
		Weight:    account.Weight,
		Email:     account.Email,
	}

	ctx.JSON(http.StatusOK, res)
}

type listAccountsReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := make([]accountResp, len(accounts)) 
	for i, v := range accounts {
		res[i] = accountResp{
			Name: v.Name,
			ID:  v.ID,
			Email: v.Email,
			Weight: v.Weight,
			BodyFat: v.BodyFat,
			StartDate: v.StartDate,
		}
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req getAccountReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, errorResponse(err))
		return
	}

	acc, err := server.store.DeleteAccount(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, acc)
}
