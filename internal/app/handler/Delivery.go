package handler

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"BMSTU_IU5_53B_rip/internal/app/models"
	"BMSTU_IU5_53B_rip/internal/app/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
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
	reqCount, _ := h.Repository.GetDeliveryReqCount(ds.DraftStatus, uint(userId))

	// Проверка на наличие заявки
	reqID, _ := h.Repository.HasRequestByUserID(uint(userId))
	// Если заявки нет, нужно вывести заявку с нулевыми полями, пустую

	var cards *[]ds.DeliveryItem
	var err error
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
	// считываем id из запроса
	id, _ := strconv.Atoi(ctx.Param("id"))
	request.ID = uint(id)

	// Привязать данные из запроса к структуре
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Проверка, что поле Image не является nil
	if request.Image == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No image in request"})
		return
	}

	// Инициализация Minio хранилища
	minioStorage, err := storage.NewMinioStorage(
		os.Getenv("MINIO_ENDPOINT_URL"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("MINIO_SECURE") == "true",
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Minio client"})
		return
	}
	fmt.Println(request.ID, "IDDDDDDDDDDDD")

	// Извлечение файла из запроса
	file, err := request.Image.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open image"})
		return
	}
	defer file.Close()

	// Генерация имени файла
	fileExtension := filepath.Ext(request.Image.Filename)
	fileName := strconv.Itoa(int(request.ID)) + fileExtension

	// Загрузка файла в Minio
	err = minioStorage.LoadImg(os.Getenv("MINIO_BUCKET_NAME"), fileName, file, request.Image.Size)
	//fmt.Println(request.Image, "Image")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load image"})
		return
	}

	// Генерация URL изображения
	imageURL := "http://" + os.Getenv("MINIO_ENDPOINT_URL") + "/" + os.Getenv("MINIO_BUCKET_NAME") + "/" + fileName

	delivery, err := h.Repository.GetDeliveryItemByID(strconv.Itoa(int(request.ID)))
	if delivery == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
		return
	}
	strURL, _ := h.Repository.UploadImage(strconv.Itoa(int(request.ID)), imageURL)

	// Ответ с URL изображения
	ctx.JSON(http.StatusOK, gin.H{"image_url": strURL})
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
