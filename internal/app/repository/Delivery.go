package repository

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"github.com/go-playground/validator/v10"
	"strconv"
)

// услуги

// DeliveryItemList возвращает список услуг
func (r *Repository) DeliveryItemList() (*[]ds.DeliveryItem, error) {
	var deliveryItems []ds.DeliveryItem
	r.db.Where("is_delete = ?", false).Order("title ASC").Find(&deliveryItems)
	return &deliveryItems, nil
}

// SearchDeliveryItem возвращает список услуг, отфильтрованный по цене
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

// DeleteDeliveryItem  удаляет услугу
func (r *Repository) DeleteDeliveryItem(id string) error {
	query := "UPDATE delivery_items SET is_delete = true WHERE id = $1"
	result := r.db.Exec(query, id)
	r.logger.Info("Rows affected:", result.RowsAffected)
	return nil
}

// GetDeliveryItemByID возвращает услугу по ID
func (r *Repository) GetDeliveryItemByID(id string) (*ds.DeliveryItem, error) {
	var DelItem ds.DeliveryItem
	intID, _ := strconv.Atoi(id)
	r.db.Find(&DelItem, intID)
	print(DelItem.ID, "ID")
	return &DelItem, nil
}

// CreateDeliveryItem создает услугу
func (r *Repository) CreateDeliveryItem(delivery *ds.DeliveryItem) (*ds.DeliveryItem, error) {
	validate := validator.New()
	err := validate.Struct(delivery)
	if err != nil {
		return nil, err
	}

	result := r.db.Create(delivery)
	if result.Error != nil {
		return nil, result.Error
	}

	return delivery, nil
}

// UploadImage загружает изображение в Minio
func (r *Repository) UploadImage(id string, img string) (string, error) {
	query := "UPDATE delivery_items SET image = $1 WHERE id = $2"
	result := r.db.Exec(query, img, id)
	r.logger.Info("Rows affected:", result.RowsAffected)

	// получить строку, которая в итоге у нас получилась
	var imageURL string
	r.db.Model(&ds.DeliveryItem{}).Where("id = ?", id).Select("image").Row().Scan(&imageURL)
	return imageURL, nil
}

// UpdateDeliveryItem обновляет услугу
func (r *Repository) UpdateDeliveryItem(delivery *ds.DeliveryItem) (*ds.DeliveryItem, error) {
	// Получаем текущий элемент доставки
	currentDeliveryItem := &ds.DeliveryItem{}
	result := r.db.First(currentDeliveryItem, delivery.ID)
	if result.Error != nil {
		return nil, result.Error
	}

	// Жизнь без костыля - не жизнь
	delivery.Image = currentDeliveryItem.Image

	validate := validator.New()
	err := validate.Struct(delivery)
	if err != nil {
		return nil, err
	}

	result = r.db.Save(delivery)
	if result.Error != nil {
		return nil, result.Error
	}

	return delivery, nil
}
