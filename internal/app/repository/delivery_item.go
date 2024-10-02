package repository

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"strconv"
)

// услуги

func (r *Repository) DeliveryItemList() (*[]ds.DeliveryItem, error) {
	var deliveryItems []ds.DeliveryItem
	r.db.Where("is_delete = ?", false).Find(&deliveryItems)
	return &deliveryItems, nil
}

func (r *Repository) SearchDeliveryItem(priceFrom, priceTo string) (*[]ds.DeliveryItem, error) {
	intPriceFrom, _ := strconv.Atoi(priceFrom)
	intPriceTo, _ := strconv.Atoi(priceTo)

	var deliveryItems []ds.DeliveryItem
	// сохраняем данные из бд в массив
	r.db.Find(&deliveryItems)

	var filteredItems []ds.DeliveryItem
	for _, item := range deliveryItems {
		if item.Price <= intPriceTo && item.Price >= intPriceFrom {
			filteredItems = append(filteredItems, item)
		}
	}
	return &filteredItems, nil
}

func (r *Repository) DeleteDeliveryItem(id string) error {
	query := "UPDATE delivery_items SET is_delete = true WHERE id = $1"
	result := r.db.Exec(query, id)
	r.logger.Info("Rows affected:", result.RowsAffected)
	return nil
}

func (r *Repository) GetDeliveryItemByID(id string) (*ds.DeliveryItem, error) {
	var DelItem ds.DeliveryItem
	intID, _ := strconv.Atoi(id)
	r.db.Find(&DelItem, intID)
	print(DelItem.ID, "ID")
	return &DelItem, nil
}

/*
func (r *Repository) CreateOrUpdateDeliveryReq(itemID uint) (*ds.DeliveryRequest, error) {
	var order ds.DeliveryRequest
	err := r.db.Where("status = ?", ds.DraftStatus).First(&order).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create a new order
		order = ds.DeliveryRequest{
			Status:      ds.DraftStatus,
			DateCreated: time.Now(),
		}
		if err := r.db.Create(&order).Error; err != nil {
			return nil, err
		}
	}

	// Add the item to the order
	itemRequest := ds.Item_request{
		ItemID:    itemID,
		RequestID: order.ID,
		Count:     1,
	}
	if err := r.db.Create(&itemRequest).Error; err != nil {
		return nil, err
	}

	return &order, nil
}*/

func (r *Repository) GetDeliveryReqLength(status string, user_id uint) (int64, error) {
	var count int64
	var req ds.DeliveryRequest

	if err := r.db.Where("user_id = ? AND status = ?", user_id, status).First(&req).Error; err != nil {
		return 0, err
	}

	reqID := req.ID

	err := r.db.Model(&ds.Item_request{}).Where("request_id = ?", reqID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
