package ds

import (
	"encoding/json"
	"time"
)

type DeliveryRequest struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	DateCreated  time.Time `json:"date_created"`
	DateFormed   time.Time `json:"date_formed"`
	DateAccepted time.Time `json:"date_accepted"`
	Status       string    `json:"status" gorm:"type:varchar(255)"`

	Address      string    `json:"address" gorm:"type:varchar(255)"`
	DeliveryDate time.Time `json:"delivery_date"`
	DeliveryType string    `json:"delivery_type" gorm:"type:varchar(255)"`
	UserID       uint      `json:"-"`
	ModeratorID  *uint     `json:"-"`
	User         User      `json:"-" gorm:"foreignKey:UserID"`
	Moderator    User      `json:"-" gorm:"foreignKey:ModeratorID"`
	DeliveryTime int64     `json:"delivery_time"`
}

const (
	DraftStatus     = "черновик"
	DeletedStatus   = "удален"
	FormedStatus    = "сформирован"
	CompletedStatus = "завершен"
	RejectedStatus  = "отклонен"
)

const (
	HomeDelivery    = "На дом"
	CourierDelivery = "Курьер"
	CarDelivery     = "Грузовик"
)

func (d DeliveryRequest) MarshalJSON() ([]byte, error) {
	type Alias DeliveryRequest
	return json.Marshal(&struct {
		DeliveryDate string `json:"delivery_date"`
		*Alias
	}{
		DeliveryDate: d.DeliveryDate.Format("2006-01-02"),
		Alias:        (*Alias)(&d),
	})
}
