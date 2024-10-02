package ds

type Item_request struct {
	ItemID    uint `json:"item_id" gorm:"primaryKey;auto_increment:false"`
	RequestID uint `json:"request_id" gorm:"primaryKey;auto_increment:false"`
	Count     int  `json:"count" gorm:"default:1"`
}
