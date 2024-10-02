package repository

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"strconv"
)

// услуги

func (r *Repository) DeliveryItemList() (*[]ds.DeliveryItem, error) {
	var deliveryItems []ds.DeliveryItem
	//r.db.Where("is_delete = ?", false).Find(&deliveryItems)
	r.db.Find(&deliveryItems)
	r.logger.Info("Found items:", len(deliveryItems))
	return &deliveryItems, nil
}

func (r *Repository) SearchDeliveryItem(price_from, price_to string) (*[]ds.DeliveryItem, error) {
	int_price_from, _ := strconv.Atoi(price_from)
	int_price_to, _ := strconv.Atoi(price_to)

	var deliveryItems []ds.DeliveryItem
	// сохраняем данные из бд в массив
	r.db.Find(&deliveryItems)

	var filteredItems []ds.DeliveryItem
	for _, item := range deliveryItems {
		if item.Price <= int_price_to && item.Price >= int_price_from {
			filteredItems = append(filteredItems, item)
		}
	}
	return &filteredItems, nil
}

func (r *Repository) DeleteDeliveryItem(id string) {
	// в данном случае не горм, а a+sql запрос - по заданию
	query := "UPDATE delivery_items SET is_delete = true WHERE id = $1"
	r.db.Exec(query, id)
}

func (r *Repository) GetDeliveryItemByID(id string) (*ds.DeliveryItem, error) {
	var DelItem ds.DeliveryItem
	intID, _ := strconv.Atoi(id)
	r.db.Find(&DelItem, intID)
	return &DelItem, nil
}
