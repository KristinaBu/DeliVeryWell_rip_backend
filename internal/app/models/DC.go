package models

type DeleteDCRequest struct {
	DeliveryID uint `json:"delivery_id"`
	CallID     uint `json:"call_id"`
}

type UpdateDCCountRequest struct {
	DeliveryID uint `json:"delivery_id"`
	CallID     uint `json:"call_id"`
	Count      int  `json:"count"`
}
