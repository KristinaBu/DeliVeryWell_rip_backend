package models

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"mime/multipart"
)

// Запрос

type GetAllDeliveryRequest struct {
	PriceFrom string `json:"price_from"`
	PriceTo   string `json:"price_to"`
}

type GetDeliveryRequest struct {
	ID string `json:"id"`
}

type CreateDeliveryRequest struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

type UploadImageRequest struct {
	ID    uint                  `json:"id"`
	Image *multipart.FileHeader `form:"image"` // Поле для файла
}

// Ответ

type GetAllDeliveryResponse struct {
	Payload *[]ds.DeliveryItem `json:"payload"`
}

type GetDeliveryResponse struct {
	Payload *ds.DeliveryItem `json:"payload"`
}

type CreateDeliveryResponse struct {
	Delivery *ds.DeliveryItem `json:"delivery"`
}

type UploadImageResponse struct {
	ImageURL string `json:"image_url"`
}

type AddDeliveryToCallResponse struct {
	DeliveryItem    *ds.DeliveryItem    `json:"delivery"`
	DeliveryRequest *ds.DeliveryRequest `json:"request"`
}
