package handler

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"BMSTU_IU5_53B_rip/internal/app/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Функция обработчика
func (h *Handler) GetAllDelivery(ctx *gin.Context) {
	var request models.GetAllDeliveryRequest
	priceFrom := ctx.Query("price_from")
	priceTo := ctx.Query("price_to")
	request.PriceFrom = priceFrom
	request.PriceTo = priceTo

	userId := 1
	reqCount, err := h.Repository.GetDeliveryReqCount(ds.DraftStatus, uint(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}
	reqID, err_ := h.Repository.HasRequestByUserID(uint(userId))
	if err_ != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err_.Error(),
		})
		return

	}

	var cards *[]ds.DeliveryItem

	if request.PriceFrom == "" && request.PriceTo == "" {
		cards, err = h.Repository.DeliveryItemList()
	} else {
		cards, err = h.Repository.SearchDeliveryItem(request.PriceFrom, request.PriceTo)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := models.GetAllDeliveryResponse{
		ReqID:        int(reqID),
		ReqCallCount: int(reqCount),
		Payload:      cards,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) GetDelivery(ctx *gin.Context) {
	var request models.GetDeliveryRequest
	request.ID = ctx.Param("id")
	card, err := h.Repository.GetDeliveryItemByID(request.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, models.GetDeliveryResponse{
		Payload: card,
	})
}

// CreateDelivery создает карточку
func (h *Handler) CreateDelivery(ctx *gin.Context) {
	var request models.CreateDeliveryRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	card := &ds.DeliveryItem{
		Title:       request.Title,
		Price:       request.Price,
		Description: request.Description,
	}
	newCard, err_ := h.Repository.CreateDeliveryItem(card)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err_.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, models.CreateDeliveryResponse{
		Delivery: newCard,
	})
}

// UploadImage загружает изображение в minio
func (h *Handler) UploadImage(ctx *gin.Context) {
	var request models.UploadImageRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Извлекаем ID из параметров маршрута
	id := ctx.Param("id")

	ImageURL, err_ := h.Repository.UploadImage(id, request.ImageURL)
	if err_ != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err_.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, models.UploadImageResponse{
		ImageURL: ImageURL,
	})
}

// UpdateDelivery обновляет карточку
func (h *Handler) UpdateDelivery(ctx *gin.Context) {
	var request models.CreateDeliveryRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	id, _ := strconv.Atoi(ctx.Param("id"))
	request.ID = uint(id)
	card := &ds.DeliveryItem{
		ID:          request.ID,
		Title:       request.Title,
		Price:       request.Price,
		Description: request.Description,
	}
	updatedCard, err_ := h.Repository.UpdateDeliveryItem(card)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err_.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, models.CreateDeliveryResponse{
		Delivery: updatedCard,
	})

}

// DeleteDelivery удаляет карточку и редиректит на главную
func (h *Handler) DeleteDelivery(ctx *gin.Context) {
	id := ctx.Param("id")
	err := h.Repository.DeleteDeliveryItem(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})
}

// AddDeliveryToCall добавляет карточку в заявку
func (h *Handler) AddDeliveryToCall(ctx *gin.Context) {
	itemID := ctx.Param("id")
	intItemID, _ := strconv.Atoi(itemID)
	userID := 1

	err := h.Repository.LinkItemToDraftRequest(uint(userID), uint(intItemID))
	if err != nil {
	}
	delivery, err_ := h.Repository.GetDeliveryItemByID(itemID)
	if err_ != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err_.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, models.AddDeliveryToCallResponse{
		DeliveryItem: delivery,
	})
}
