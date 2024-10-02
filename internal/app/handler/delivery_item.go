package handler

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"github.com/gin-gonic/gin"
	"net/http"
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

/*
	func (h *Handler) GetMyCallCards(ctx *gin.Context) {
		if callRequestId, err := strconv.Atoi(ctx.Param("callrequest_id")); err == nil {
			callRequest := h.Repository.GetMyCallCards(callRequestId)
			ctx.HTML(http.StatusOK, "mycards.html", gin.H{
				"payload":      callRequest.Cards,
				"Data":         callRequest.Data,
				"Address":      callRequest.Address,
				"DeliveryName": callRequest.DeliveryName,
			})
		} else {
			h.errorHandler(ctx, http.StatusBadRequest, err)
		}
	}
*/

/*

func (h *Handler) CreateOrUpdateDeliveryReq(ctx *gin.Context) {
	var request struct {
		ItemID uint `json:"item_id"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.Repository.CreateOrUpdateDeliveryReq(request.ItemID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
}
*/
