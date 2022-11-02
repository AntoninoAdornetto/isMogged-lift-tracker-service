package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ParseUUIDStr(uuidStr string, ctx *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(uuidStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	return id, err
}
