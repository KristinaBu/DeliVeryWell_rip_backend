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

// DeleteDeliveryReq удаляет заявку и редиректит на главную
func (h *Handler) DeleteDeliveryReq(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.Repository.DeleteDeliveryReq(id)
	if err != nil {
	}
	ctx.Redirect(http.StatusFound, "/")

}

// GetMyCallCards рисует страницу с заявкой
func (h *Handler) GetMyCallCards1(ctx *gin.Context) {
	if callRequestId, err := strconv.Atoi(ctx.Param("callrequest_id")); err == nil {
		print("id req = ", callRequestId)

		// Предполагаем, что пользователь идентификатор равен 1
		user_id := 1

		// Получаем заявку по ID
		callRequest, err := h.Repository.GetCallRequestById(uint(callRequestId))
		if err != nil || callRequest.Status == ds.DeletedStatus {
			// Если заявка не найдена или удалена, перенаправляем на главную страницу
			ctx.Redirect(http.StatusFound, "/")
			return
		}

		// Получаем карточки доставки для этой заявки
		cards, err := h.Repository.GetDeliveryItemsByUserAndStatus(ds.DraftStatus, uint(user_id))
		if err != nil {
			// Если произошла ошибка, перенаправляем на главную страницу
			ctx.Redirect(http.StatusFound, "/")
			return
		}

		timestamp := callRequest.DeliveryDate
		formattedTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second())

		// получаем колическво карточек из м-м таблицы item_request
		count, err := h.Repository.GetDeliveryReqCount(ds.DraftStatus, uint(user_id))

		ctx.HTML(http.StatusOK, "mycards.html", gin.H{
			"payload":      cards,
			"Data":         formattedTime,
			"Address":      callRequest.Address,
			"DeliveryType": callRequest.DeliveryType,
			"ReqID":        callRequestId,
			"Count":        count,
			"CallDomain":   CallDomain,
		})
	} else {
		h.errorHandler(ctx, http.StatusBadRequest, err)
	}
}

// GetMyCallCards рисует страницу с заявкой
func (h *Handler) GetMyCallCards(ctx *gin.Context) {}
