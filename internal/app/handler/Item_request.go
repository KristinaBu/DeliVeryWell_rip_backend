package handler

import (
	"BMSTU_IU5_53B_rip/internal/app/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Ping godoc
// @Summary Delete delivery call
// @Description delete delivery call

func (h *Handler) DeleteDC(ctx *gin.Context) {
	callId, _ := strconv.Atoi(ctx.Param("id"))
	var request models.DeleteDCRequest
	request.CallID = uint(callId)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Repository.DeleteDC(request.DeliveryID, request.CallID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// UpdateDCCount - обновляет количество услуг в заявке
func (h *Handler) UpdateDCCount(ctx *gin.Context) {
	callId, _ := strconv.Atoi(ctx.Param("id"))
	var request models.UpdateDCCountRequest
	request.CallID = uint(callId)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Repository.UpdateDCCount(request.DeliveryID, request.CallID, request.Count)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
