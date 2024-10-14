package ds

type Item_request struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	ItemID    uint `json:"item_id" gorm:"uniqueIndex:item_request_key"`
	RequestID uint `json:"request_id" gorm:"uniqueIndex:item_request_key"`
	Count     int  `json:"count" gorm:"default:1"`

	Item    DeliveryItem    `json:"-" gorm:"foreignKey:ItemID"`
	Request DeliveryRequest `json:"-" gorm:"foreignKey:RequestID"`
}
