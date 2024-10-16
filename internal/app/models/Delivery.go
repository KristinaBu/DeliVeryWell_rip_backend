package models

import "BMSTU_IU5_53B_rip/internal/app/ds"

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
	ImageURL string `json:"image_url"`
}

// Ответ

type GetAllDeliveryResponse struct {
	ReqID        int                `json:"req_id"`
	ReqCallCount int                `json:"req_call_count"`
	Payload      *[]ds.DeliveryItem `json:"payload"`
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
	DeliveryItem *ds.DeliveryItem `json:"delivery"`
}
