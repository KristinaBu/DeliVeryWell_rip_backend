package models

import "BMSTU_IU5_53B_rip/internal/app/ds"

// Запросы

type GetCallsRequest struct {
	DateFrom string `form:"date_from"` // дата начала диапазона
	DateTo   string `form:"date_to"`   // дата конца диапазона
	Status   string `form:"status"`    // статус
}

type GetMyCallCardsRequest struct {
	UserId string `form:"user_id"` // идентификатор пользователя
}

type UpdateCallRequest struct {
	ID           uint   `json:"id"`            // идентификатор
	Address      string `json:"address"`       // адрес
	DeliveryDate string `json:"delivery_date"` // дата доставки
	DeliveryType string `json:"delivery_type"` // тип доставки
}

type FinishCallRequest struct {
	ID     uint `json:"id"`      // идентификатор
	UserID uint `json:"user_id"` // модератор
}

type CompleteOrRejectCallRequest struct {
	ID          uint `json:"id"`           // идентификатор
	ModeratorID uint `json:"moderator_id"` // модератор
	IsComplete  bool `json:"is_complete"`  // завершен
}

// Ответы

type GetCallsResponse struct {
	Calls []*ds.DeliveryRequest `json:"calls"` // список звонков
}

type GetMyCallCardsResponse struct {
	CallRequest   *ds.DeliveryRequest `json:"call_request"`   // заявка на доставку
	DeliveryItems []*ds.DeliveryItem  `json:"delivery_items"` // карточки доставки
	Count         int                 `json:"count"`          // количество карточек
}

type GetCallResponse struct {
	CallRequest   *ds.DeliveryRequest     `json:"call_request"`   // заявка на доставку
	DeliveryItems []DeliveryItemWithCount `json:"delivery_items"` // карточки доставки
	// общее число доставок
	DeliveriesCount int `json:"deliveries_count"`
}

type DeliveryItemWithCount struct {
	ds.DeliveryItem
	Count int `json:"count"`
}

type UpdateCallResponse struct {
	CallRequest *ds.DeliveryRequest `json:"call_request"` // заявка на доставку
}

type FinishCallResponse struct {
	CallRequest *ds.DeliveryRequest `json:"call_request"` // заявка на доставку
}

type CompleteOrRejectCallResponse struct {
	CallRequest *ds.DeliveryRequest `json:"call_request"` // заявка на доставку
	TotalCount  int                 `json:"total_count"`  // общее количество
}
