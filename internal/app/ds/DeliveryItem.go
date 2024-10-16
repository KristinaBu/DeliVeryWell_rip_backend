package ds

type DeliveryItem struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Image       string `json:"image" gorm:"type:varchar(255);default:''"`
	Title       string `json:"title" gorm:"type:varchar(255)"`
	Price       int    `json:"price" gorm:"not null"`
	Description string `json:"description" gorm:"type:text"`
	IsDelete    bool   `json:"is_delete" gorm:"default:false"`
}
