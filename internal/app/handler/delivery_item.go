package handler

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// вызываются функции из репы, которые идут в бд
// то есть как бы прослойка между эндпоинтами и данными, которые идут их бд

func (h *Handler) DeliveryItemList(ctx *gin.Context) {
	priceFrom := ctx.Query("price_from")
	priceTo := ctx.Query("price_to")

	user_id := 1
	reqCount, _ := h.Repository.GetDeliveryReqLength(ds.DraftStatus, uint(user_id))

	if priceFrom == "" && priceTo == "" {
		cards, err := h.Repository.DeliveryItemList()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"NoCards":      "",
			"payload":      cards,
			"SearchFrom":   priceFrom,
			"SearchUp":     priceTo,
			"ReqCallCount": reqCount,
		})
		return
	}
	cards, err := h.Repository.SearchDeliveryItem(priceFrom, priceTo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"NoCards":      "",
		"payload":      cards,
		"SearchFrom":   priceFrom,
		"SearchUp":     priceTo,
		"ReqCallCount": reqCount,
	})

}

func (h *Handler) DeliveryItemByID(ctx *gin.Context) {
	id := ctx.Param("id")
	card, err := h.Repository.GetDeliveryItemByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	print("card", card.ID, card.Title)
	ctx.HTML(http.StatusOK, "cardDetails.html", gin.H{
		"payload": card,
	})
}

func (h *Handler) DeleteDeliveryItem(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.Repository.DeleteDeliveryItem(id)
	if err != nil {
	}
	ctx.Redirect(http.StatusFound, "/")
}

func (h *Handler) DeleteDeliveryReq(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.Repository.DeleteDeliveryReq(id)
	if err != nil {
	}
	fmt.Println("1ID del req   ", id, " stetus ")
	ctx.Redirect(http.StatusFound, "/")

}

func (h *Handler) AddDeliveryItem(ctx *gin.Context) {
	itemID := ctx.Param("id")
	intItemID, _ := strconv.Atoi(itemID)
	userID := 1

	err := h.Repository.LinkItemToDraftRequest(uint(userID), uint(intItemID))
	if err != nil {
	}
	fmt.Println("1ID del req   ", itemID, " stetus ")
	ctx.Redirect(http.StatusFound, "/")

}

func (h *Handler) GetMyCallCards(ctx *gin.Context) {
	if callRequestId, err := strconv.Atoi(ctx.Param("callrequest_id")); err == nil {
		print("id req = ", callRequestId)

		// Предполагаем, что пользователь идентификатор равен 1
		user_id := 1

		// Получаем заявку по ID
		callRequest, err := h.Repository.GetCallRequestById(uint(callRequestId))
		if err != nil {
			h.errorHandler(ctx, http.StatusNotFound, err)
			return
		}

		// Получаем карточки доставки для этой заявки
		cards, err := h.Repository.GetDeliveryItemsByUserAndStatus(ds.DraftStatus, uint(user_id))
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.HTML(http.StatusOK, "mycards.html", gin.H{
			"payload":      cards,
			"Data":         callRequest.DeliveryDate,
			"Address":      callRequest.Address,
			"DeliveryType": callRequest.DeliveryType,
			"ReqID":        callRequestId,
			"Count":        1,
		})
	} else {
		h.errorHandler(ctx, http.StatusBadRequest, err)
	}
}
