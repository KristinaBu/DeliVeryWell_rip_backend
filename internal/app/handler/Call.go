package handler

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"BMSTU_IU5_53B_rip/internal/app/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetCalls - возращает все заявки
func (h *Handler) GetCalls(ctx *gin.Context) {
	var request models.GetCallsRequest
	dateFromQuery := ctx.Query("date_from")
	dateToQuery := ctx.Query("date_to")
	statusQuery := ctx.Query("status")

	request.DateFrom = dateFromQuery
	request.DateTo = dateToQuery
	request.Status = statusQuery
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"
	dateFrom, err := time.Parse(layout, request.DateFrom)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dateTo, err := time.Parse(layout, request.DateTo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	calls, err := h.Repository.GetCalls(dateFrom, dateTo, request.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.GetCallsResponse{Calls: calls})
}

// DeleteCall - устанавливает статус "удалено" для звонка
func (h *Handler) DeleteCall(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.Repository.DeleteCall(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetMyCallCards рисует страницу с заявкой
func (h *Handler) GetMyCallCards(ctx *gin.Context) {
	if callRequestId, err := strconv.Atoi(ctx.Param("id")); err == nil {
		// Предполагаем, что пользователь идентификатор равен 1
		userId := 1

		// Получаем заявку по ID
		callRequest, err := h.Repository.GetCallRequestById(uint(callRequestId))
		if err != nil || callRequest.Status == ds.DeletedStatus {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Call request not found or deleted"})
			return
		}

		// Получаем карточки доставки для этой заявки
		cards, err := h.Repository.GetDeliveryItemsByUserAndStatus(ds.DraftStatus, uint(userId))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// получаем колическво карточек из м-м таблицы item_request
		count, err := h.Repository.GetDeliveryReqCount(ds.DraftStatus, uint(userId))

		ctx.JSON(http.StatusOK, models.GetMyCallCardsResponse{
			CallRequest:   callRequest,
			DeliveryItems: cards,
			Count:         int(count),
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// GetCall возвращает заявку на звонок
func (h *Handler) GetCall(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	call, err := h.Repository.GetCallRequestById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	deliveries, err := h.Repository.GetDeliveryItemsByCallRequestID(uint(id))
	ctx.JSON(http.StatusOK, models.GetMyCallCardsResponse{
		CallRequest:   call,
		DeliveryItems: deliveries,
		Count:         len(deliveries),
	})
}

// UpdateCall обновляет заявку на звонок по теме
func (h *Handler) UpdateCall(ctx *gin.Context) {
	var request models.UpdateCallRequest
	id, _ := strconv.Atoi(ctx.Param("id"))
	request.ID = uint(id)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("даты:", request.DeliveryDate)
	layout := "2006-01-02" // Измените формат на этот, если вы передаете дату без времени
	deliveryDate, err := time.Parse(layout, request.DeliveryDate)
	if err != nil {
		fmt.Println("Ошибка при парсинге даты:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	call := &ds.DeliveryRequest{
		ID:           request.ID,
		Address:      request.Address,
		DeliveryDate: deliveryDate,
		DeliveryType: request.DeliveryType,
	}
	resp, err := h.Repository.UpdateCall(call)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.UpdateCallResponse{
		CallRequest: resp,
	})
}

// FormCall - формирует заявку на звонок
func (h *Handler) FormCall(ctx *gin.Context) {
	var request models.FinishCallRequest
	id, _ := strconv.Atoi(ctx.Param("id"))
	request.ID = uint(id)

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(request.UserID, "")
	resp, err := h.Repository.FormCall(request.ID, request.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.UpdateCallResponse{
		CallRequest: resp,
	})
}

// CompleteOrRejectCall - завершает заявку на звонок
func (h *Handler) CompleteOrRejectCall(ctx *gin.Context) {
	var request models.CompleteOrRejectCallRequest
	id, _ := strconv.Atoi(ctx.Param("id"))
	request.ID = uint(id)
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем полный объект call из базы данных
	call, err := h.Repository.GetCallRequestById(request.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем ModeratorID из запроса
	call.ModeratorID = request.ModeratorID

	resp, totalCount, err := h.Repository.CompleteOrRejectCall(call, request.IsComplete)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.CompleteOrRejectCallResponse{
		CallRequest: resp,
		TotalCount:  totalCount,
	})
}
