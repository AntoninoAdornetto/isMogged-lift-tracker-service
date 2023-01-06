package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/AntoninoAdornetto/isMogged-lift-tracker-service/db/sqlc"
	"github.com/AntoninoAdornetto/isMogged-lift-tracker-service/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRes struct {
	ID                uuid.UUID `json:"user_id"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	StartDate         time.Time `json:"start_date"`
}

func newUserResponse(account db.GetAccountByEmailRow) UserRes {
	return UserRes{
		ID:                account.ID,
		Email:             account.Email,
		PasswordChangedAt: account.PasswordChangedAt,
		StartDate:         account.StartDate,
	}
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRes struct {
	JWT  string  `json:"jwt"`
	User UserRes `json:"user"`
}

func (server *Server) login(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccountByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.ValidatePassword(req.Password, account.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	jwt, err := server.tokenCreator.CreateToken(account.ID, server.config.AccessDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := LoginRes{
		JWT:  jwt,
		User: newUserResponse(account),
	}

	ctx.JSON(http.StatusOK, res)
}
