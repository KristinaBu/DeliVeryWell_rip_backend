package models

import "BMSTU_IU5_53B_rip/internal/app/ds"

// Запросы

type GetCallsRequest struct {
	DateFrom string `form:"date_from"` // дата начала диапазона
	DateTo   string `form:"date_to"`   // дата конца диапазона
	Status   string `form:"status"`    // статус
}

// Ответы

type GetCallsResponse struct {
	Calls []*ds.DeliveryRequest `json:"calls"` // список звонков
}
