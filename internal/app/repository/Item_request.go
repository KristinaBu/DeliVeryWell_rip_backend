package repository

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"fmt"
)

func (r *Repository) DeleteDC(deliveryID uint, callID uint) error {
	// Получение услуги из базы данных
	var existingDC ds.Item_request
	if err := r.db.Where("item_id = ? AND request_id = ?", deliveryID, callID).First(&existingDC).Error; err != nil {
		return err
	}

	// Получение заявки из базы данных
	var existingRequest ds.DeliveryRequest
	if err := r.db.First(&existingRequest, callID).Error; err != nil {
		return err
	}

	// Проверка, что статус заявки является "сформирован" или "черновик"
	if existingRequest.Status != ds.FormedStatus && existingRequest.Status != ds.DraftStatus {
		return fmt.Errorf("request status is not formed or draft")
	}

	// Удаление услуги из базы данных
	result := r.db.Delete(&existingDC)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) UpdateDCCount(deliveryID uint, callID uint, count int) error {
	if count <= 0 {
		return fmt.Errorf("not positive count")
	}
	// Получение услуги из базы данных
	var existingDC ds.Item_request
	if err := r.db.Where("item_id = ? AND request_id = ?", deliveryID, callID).First(&existingDC).Error; err != nil {
		return err
	}

	// Получение заявки из базы данных
	var existingRequest ds.DeliveryRequest
	if err := r.db.First(&existingRequest, callID).Error; err != nil {
		return err
	}

	// Проверка, что статус заявки является "сформирован" или "черновик"
	if existingRequest.Status != ds.FormedStatus && existingRequest.Status != ds.DraftStatus {
		return fmt.Errorf("request status is not formed or draft")
	}

	// Обновление количества услуг в заявке
	result := r.db.Model(&existingDC).Update("Count", count)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateDC - создание связи между услугой и заявкой, count = 1 default
func (r *Repository) CreateDC(itemID, requestID uint) error {
	itemRequest := ds.Item_request{
		ItemID:    itemID,
		RequestID: requestID,
		Count:     1,
	}
	_ = r.db.Create(&itemRequest).Error

	return nil
}