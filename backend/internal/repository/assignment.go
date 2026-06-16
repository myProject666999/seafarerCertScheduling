package repository

import (
	"seafarer-cert-scheduling/internal/model"
)

type AssignmentRepo struct{}

func (r *AssignmentRepo) List(page, pageSize int, shipID, seafarerID int64, status int8) ([]model.SeafarerAssignment, int64, error) {
	var list []model.SeafarerAssignment
	var total int64
	db := model.DB.Model(&model.SeafarerAssignment{})
	if shipID > 0 {
		db = db.Where("ship_id = ?", shipID)
	}
	if seafarerID > 0 {
		db = db.Where("seafarer_id = ?", seafarerID)
	}
	if status >= 0 {
		db = db.Where("status = ?", status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes(paginate(page, pageSize)).Preload("Seafarer").Preload("Ship").Preload("ShipPosition").Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *AssignmentRepo) GetByID(id int64) (*model.SeafarerAssignment, error) {
	var a model.SeafarerAssignment
	if err := model.DB.Preload("Seafarer").Preload("Ship").Preload("ShipPosition").First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AssignmentRepo) Create(a *model.SeafarerAssignment) error {
	return model.DB.Create(a).Error
}

func (r *AssignmentRepo) Update(a *model.SeafarerAssignment) error {
	return model.DB.Save(a).Error
}

func (r *AssignmentRepo) FindActiveBySeafarer(seafarerID int64) (*model.SeafarerAssignment, error) {
	var a model.SeafarerAssignment
	if err := model.DB.Where("seafarer_id = ? AND status = ?", seafarerID, model.AssignmentStatusOnboard).
		Preload("Ship").Preload("ShipPosition").First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AssignmentRepo) FindActiveByPosition(positionID int64) ([]model.SeafarerAssignment, error) {
	var list []model.SeafarerAssignment
	if err := model.DB.Where("ship_position_id = ? AND status = ?", positionID, model.AssignmentStatusOnboard).
		Preload("Seafarer").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
