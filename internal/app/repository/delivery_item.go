package repository

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"time"
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

func (r *Repository) DeleteDeliveryReq(id string) error {
	query := "UPDATE delivery_requests SET status = 'удален' WHERE id = $1"
	result := r.db.Exec(query, id)
	fmt.Println("ID del req   ", id, " stetus ")
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

func (r *Repository) HasRequestByUserID(userID uint) (uint, error) {
	var req ds.DeliveryRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&req).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	return req.ID, nil
}

func (r *Repository) CreateOrUpdateDeliveryReq(itemID, userID uint) (*ds.DeliveryRequest, error) {
	var order ds.DeliveryRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&order).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// создать
		order = ds.DeliveryRequest{
			UserID:      userID,
			Status:      ds.DraftStatus,
			DateCreated: time.Now(),
		}
		if err := r.db.Create(&order).Error; err != nil {
			return nil, err
		}
	}

	// добавить в заявку
	itemRequest := ds.Item_request{
		ItemID:    itemID,
		RequestID: order.ID,
		Count:     1,
	}
	if err := r.db.Create(&itemRequest).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

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

func (r *Repository) GetDeliveryItemsByUserAndStatus(status string, userID uint) ([]*ds.DeliveryItem, error) {
	var items []*ds.DeliveryItem

	// Используем GORM для выполнения запроса
	err := r.db.Model(&ds.DeliveryItem{}).
		Select("delivery_items.*").
		Joins("INNER JOIN item_requests ON delivery_items.id = item_requests.item_id").
		Joins("INNER JOIN delivery_requests ON item_requests.request_id = delivery_requests.id").
		Where("delivery_requests.user_id = ?", userID).
		Where("delivery_requests.status = ?", status).
		Find(&items).Error

	if err != nil {
		return nil, err
	}
	fmt.Println("Len items  ", len(items))

	return items, nil
}

func (r *Repository) GetCallRequestById(id uint) (*ds.DeliveryRequest, error) {
	var callRequest ds.DeliveryRequest
	err := r.db.Where("id = ?", id).First(&callRequest).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching call request: %w", err)
	}

	// Выводим информацию о найденной записи
	r.logger.Infof("Found call request ID: %d", id)

	return &callRequest, nil
}

func (r *Repository) CreateDraftRequestAndGetID(userID uint) (uint, error) {
	draftRequest := ds.DeliveryRequest{
		Status:       ds.DraftStatus,
		UserID:       userID,
		Address:      "",
		DeliveryDate: time.Now(),
		DeliveryType: "Курьер",
	}

	err := r.db.Create(&draftRequest).Error
	if err != nil {
		return 0, fmt.Errorf("error creating draft request: %w", err)
	}

	r.logger.Infof("Created new draft request ID: %d", draftRequest.ID)

	return draftRequest.ID, nil
}

func (r *Repository) LinkItemToDraftRequest(userID uint, itemId uint) error {
	// поик существующей заявки пользователя со статусом 'черновик'
	var draftRequest ds.DeliveryRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&draftRequest).Error
	if err == gorm.ErrRecordNotFound {
		// если заявки нет, создаем новую
		draftRequest.UserID = userID
		draftRequest.Status = ds.DraftStatus
		draftRequest.Address = ""
		draftRequest.DeliveryDate = time.Now()
		draftRequest.DeliveryType = "Курьер"
		err = r.db.Create(&draftRequest).Error
		if err != nil {
			return fmt.Errorf("error creating new draft request: %w", err)
		}
		r.logger.Infof("Created new draft request ID: %d for user ID: %d", draftRequest.ID, userID)
	} else {
		r.logger.Infof("Found existing draft request ID: %d for user ID: %d", draftRequest.ID, userID)
	}

	// Добавляем элемент в существующую заявку
	itemRequest := ds.Item_request{
		ItemID:    itemId,
		RequestID: draftRequest.ID,
		Count:     1,
	}
	err = r.db.Create(&itemRequest).Error
	if err != nil {
		return fmt.Errorf("error linking item to draft request: %w", err)
	}

	return nil
}
