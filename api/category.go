package api

import (
	"context"
	"net/http"
	"strconv"

	db "github.com/AntoninoAdornetto/lift_tracker/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createCategoryReq struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) createCategory(ctx *gin.Context) {
	var req createCategoryReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.CreateCategory(context.Background(), req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type getCategoryReq struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) getCategory(ctx *gin.Context) {
	var req getCategoryReq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := strconv.Atoi(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	query, err := server.store.GetCategory(context.Background(), int16(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, query)
}

func (server *Server) listCategories(ctx *gin.Context) {
	categories, err := server.store.ListCategories(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func (server *Server) updateCategory(ctx *gin.Context) {
	var uri getCategoryReq
	var req createCategoryReq
	if err := ctx.BindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := strconv.Atoi(uri.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateCategoryParams{
		Name: req.Name,
		ID:   int16(id),
	}

	err = server.store.UpdateCategory(context.Background(), args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (server *Server) deleteCategory(ctx *gin.Context) {
	var req getCategoryReq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := strconv.Atoi(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteCategory(context.Background(), int16(id))

	ctx.JSON(http.StatusNoContent, nil)
}
