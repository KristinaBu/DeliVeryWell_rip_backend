package repository

import (
	"BMSTU_IU5_53B_rip/internal/app/ds"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// DeleteDeliveryReq  удаляет заявку
func (r *Repository) DeleteCall(id string) error {
	query := "UPDATE delivery_requests SET status = 'удален' WHERE id = $1"
	result := r.db.Exec(query, id)
	fmt.Println("ID del req   ", id, " stetus ")
	r.logger.Info("Rows affected:", result.RowsAffected)

	return nil
}

// GetDeliveryReqCount возвращает количество элементов в заявке
func (r *Repository) GetDeliveryReqCount(status string, userId uint) (int64, error) {
	var count int64
	var req ds.DeliveryRequest

	if err := r.db.Where("user_id = ? AND status = ?", userId, status).First(&req).Error; err != nil {
		return 0, err
	}

	reqID := req.ID

	err := r.db.Model(&ds.Item_request{}).Where("request_id = ?", reqID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetDeliveryItemsByUserAndStatus возвращает элементы заявки пользователя по статусу
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

// GetCallRequestById возвращает заявку по ID
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

// CreateDraftRequestAndGetID создает черновик заявки и возвращает ID
func (r *Repository) CreateDraftRequestAndGetID(userID, itemID uint) (uint, error) {
	draftRequest := ds.DeliveryRequest{
		Status:       ds.DraftStatus,
		UserID:       userID,
		Address:      "",
		DateCreated:  time.Now(),
		DeliveryType: ds.CourierDelivery,
	}
	err := r.db.Create(&draftRequest).Error
	if err != nil {
		return 0, fmt.Errorf("error creating draft request: %w", err)
	}
	r.logger.Infof("Created new draft request ID: %d", draftRequest.ID)
	// создаем связь с таблицей м-м, чтобы добавить доставку в звонок
	err = r.CreateDC(itemID, draftRequest.ID)
	if err != nil {
		return 0, err
	}

	return draftRequest.ID, nil
}

// LinkItemToDraftRequest связывает элемент с черновиком заявки, возвращаем заявку
func (r *Repository) LinkItemToDraftRequest(userID uint, itemId uint) (*ds.DeliveryRequest, error) {
	// нужно проверить, является ли доставка удаленной
	var item ds.DeliveryItem
	err_ := r.db.Where("id = ?", itemId).First(&item).Error
	if err_ != nil {
		return nil, fmt.Errorf("error fetching item: %w", err_)
	}
	if item.IsDelete == true {
		return nil, fmt.Errorf("item with id %d is deleted", itemId)
	}

	// поик существующей заявки пользователя со статусом 'черновик'
	var draftRequest ds.DeliveryRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&draftRequest).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// если заявки нет, создаем новую
		id, err := r.CreateDraftRequestAndGetID(userID, itemId)
		if err != nil {
			return nil, fmt.Errorf("error creating new draft request: %w", err)
		}
		draftRequest.ID = id
		r.logger.Infof("Created new draft request ID: %d for user ID: %d", id, userID)
	} else {
		r.logger.Infof("Found existing draft request ID: %d for user ID: %d", draftRequest.ID, userID)
	}

	err = r.CreateDC(itemId, draftRequest.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating DC: %w", err)
	}

	return &draftRequest, nil
}

// HasRequestByUserID проверяет наличие заявки пользователя
func (r *Repository) HasRequestByUserID(userID uint) (uint, error) {
	var req ds.DeliveryRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).First(&req).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Если ошибка о том, то записи нет, то нет ошибки, тк нужно потом вывести заявку с null+0 полями
			return 0, nil
		}
		return 0, err
	}
	return req.ID, nil
}

// GetMyCalls возвращает звонки пользователя в статусе черновик
func (r *Repository) GetMyCalls(userID uint) ([]*ds.DeliveryRequest, error) {
	var calls []*ds.DeliveryRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, ds.DraftStatus).Find(&calls).Error
	if err != nil {
		return nil, err
	}
	return calls, nil
}

// GetCalls - возвращает звонки с учетом фильтров
func (r *Repository) GetCalls(dateFrom, dateTo time.Time, status string, userID uint) ([]*ds.DeliveryRequest, error) {
	var calls []*ds.DeliveryRequest
	if r.IsAdmin(userID) == false {
		query := "SELECT * FROM delivery_requests WHERE user_id = ? AND status = ?"
		result := r.db.Raw(query, userID, ds.DraftStatus).Scan(&calls)
		if result.Error != nil {
			return nil, result.Error
		}
		return calls, nil
	}
	if status == "" {
		return nil, fmt.Errorf("status is empty")
	}
	query := "SELECT * FROM delivery_requests WHERE date_formed BETWEEN ? AND ? AND status = ?"
	result := r.db.Raw(query, dateFrom, dateTo, status).Scan(&calls)
	if result.Error != nil {
		return nil, result.Error
	}
	return calls, nil
}

// UpdateCall обновляет звонок
func (r *Repository) UpdateCall(call *ds.DeliveryRequest) (*ds.DeliveryRequest, error) {
	result := r.db.Model(&ds.DeliveryRequest{}).Where("id = ?", call.ID).Updates(map[string]interface{}{
		"Address":      call.Address,
		"DeliveryDate": call.DeliveryDate,
		"DeliveryType": call.DeliveryType,
	})
	if result.Error != nil {
		return nil, result.Error
	}

	// Загрузка обновленной записи из базы данных
	var updatedCall ds.DeliveryRequest
	if err := r.db.First(&updatedCall, call.ID).Error; err != nil {
		return nil, err
	}

	return &updatedCall, nil
}

// FormCall - формирование звонка
func (r *Repository) FormCall(callID uint, userID uint) (*ds.DeliveryRequest, error) {
	// Получение звонка из базы данных
	var existingCall ds.DeliveryRequest
	if err := r.db.Where("id = ? AND user_id = ?", callID, userID).First(&existingCall).Error; err != nil {
		return nil, err
	}
	fmt.Println("call", existingCall.ID, existingCall.UserID, existingCall.Status)

	b := !(existingCall.UserID == userID) || !(r.IsAdmin(userID))
	a := !(existingCall.UserID == userID)
	c := r.IsAdmin(userID)
	fmt.Println(b, a, c, existingCall.UserID, userID)
	// Проверка, что пользователь является владельцем заявки или модератором
	if !(existingCall.UserID == userID) {
		if !(r.IsAdmin(userID)) {
			return nil, fmt.Errorf("user with id %d is not the owner of the call request or a moderator", userID)
		}
	}

	// Проверка, что статус звонка является "черновиком"
	if existingCall.Status != ds.DraftStatus {
		return nil, fmt.Errorf("call status is not draft")
	}

	// рандомное число от 1 до 24
	deliveryTime := rand.Int63n(24) + 1

	result := r.db.Model(&existingCall).Updates(map[string]interface{}{
		"Status":       ds.FormedStatus,
		"DateFormed":   time.Now(),
		"DeliveryTime": deliveryTime,
	})
	if result.Error != nil {
		return nil, result.Error
	}

	// Загрузка обновленной записи из базы данных
	var updatedCall ds.DeliveryRequest
	if err := r.db.First(&updatedCall, callID).Error; err != nil {
		return nil, err
	}

	return &updatedCall, nil
}

// CompleteOrRejectCall - завершает или отклоняет заявку
func (r *Repository) CompleteOrRejectCall(call *ds.DeliveryRequest, isComplete bool) (*ds.DeliveryRequest, int, error) {
	// Проверка, является ли пользователь администратором
	isAdmin := r.IsAdmin(*call.ModeratorID)

	if !isAdmin {
		return nil, 0, fmt.Errorf("user with id %d is not an admin", call.ModeratorID)
	}

	// Получение звонка из базы данных
	var existingCall ds.DeliveryRequest
	if err := r.db.First(&existingCall, call.ID).Error; err != nil {
		return nil, 0, err
	}

	// Проверка, что статус звонка является "сформированным"
	if existingCall.Status != ds.FormedStatus {
		return nil, 0, fmt.Errorf("call status is not formed")
	}

	// Обновление статуса звонка
	newStatus := ds.RejectedStatus
	if isComplete {
		newStatus = ds.CompletedStatus
	}

	result := r.db.Model(&existingCall).Updates(map[string]interface{}{
		"ModeratorID":  call.ModeratorID,
		"Status":       newStatus,
		"DateAccepted": time.Now(),
	})
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Вычисление итогового количества единиц доставки
	var totalItemCount int
	r.db.Model(&ds.Item_request{}).Where("request_id = ?", call.ID).Select("sum(count)").Row().Scan(&totalItemCount)

	// Загрузка обновленной записи из базы данных
	var updatedCall ds.DeliveryRequest
	if err := r.db.First(&updatedCall, call.ID).Error; err != nil {
		return nil, 0, err
	}

	return &updatedCall, totalItemCount, nil
}

// GetItemRequestsByCallRequestID возвращает элементы заявки по ID звонка-заявки
func (r *Repository) GetItemRequestsByCallRequestID(callRequestID uint) ([]ds.Item_request, error) {
	var itemRequests []ds.Item_request
	err := r.db.Where("request_id = ?", callRequestID).Preload("Item").Find(&itemRequests).Error
	if err != nil {
		return nil, err
	}
	return itemRequests, nil
}
