package repository

import (
	"seafarer-cert-scheduling/internal/model"
)

type ContractRepo struct{}

func (r *ContractRepo) List(page, pageSize int, seafarerID, shipID int64) ([]model.VoyageContract, int64, error) {
	var list []model.VoyageContract
	var total int64
	db := model.DB.Model(&model.VoyageContract{})
	if seafarerID > 0 {
		db = db.Where("seafarer_id = ?", seafarerID)
	}
	if shipID > 0 {
		db = db.Where("ship_id = ?", shipID)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes(paginate(page, pageSize)).Preload("Seafarer").Preload("Ship").Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *ContractRepo) Create(c *model.VoyageContract) error {
	return model.DB.Create(c).Error
}

func (r *ContractRepo) GetByID(id int64) (*model.VoyageContract, error) {
	var c model.VoyageContract
	if err := model.DB.Preload("Seafarer").Preload("Ship").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ContractRepo) Update(c *model.VoyageContract) error {
	return model.DB.Save(c).Error
}

type EmbarkRecordRepo struct{}

func (r *EmbarkRecordRepo) List(seafarerID int64) ([]model.EmbarkDisembarkRecord, error) {
	var list []model.EmbarkDisembarkRecord
	db := model.DB.Order("record_date DESC, id DESC")
	if seafarerID > 0 {
		db = db.Where("seafarer_id = ?", seafarerID)
	}
	if err := db.Preload("Seafarer").Preload("Ship").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *EmbarkRecordRepo) Create(rec *model.EmbarkDisembarkRecord) error {
	return model.DB.Create(rec).Error
}

type LeaveRepo struct{}

func (r *LeaveRepo) List(seafarerID int64) ([]model.LeaveRecord, error) {
	var list []model.LeaveRecord
	db := model.DB.Order("id DESC")
	if seafarerID > 0 {
		db = db.Where("seafarer_id = ?", seafarerID)
	}
	if err := db.Preload("Seafarer").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *LeaveRepo) GetActiveLeave(seafarerID int64) (*model.LeaveRecord, error) {
	var lr model.LeaveRecord
	if err := model.DB.Where("seafarer_id = ? AND status = ?", seafarerID, model.LeaveStatusActive).First(&lr).Error; err != nil {
		return nil, err
	}
	return &lr, nil
}

func (r *LeaveRepo) Create(rec *model.LeaveRecord) error {
	return model.DB.Create(rec).Error
}

func (r *LeaveRepo) Update(rec *model.LeaveRecord) error {
	return model.DB.Save(rec).Error
}

type HealthRepo struct{}

func (r *HealthRepo) List(seafarerID int64) ([]model.HealthReexamination, error) {
	var list []model.HealthReexamination
	db := model.DB.Order("id DESC")
	if seafarerID > 0 {
		db = db.Where("seafarer_id = ?", seafarerID)
	}
	if err := db.Preload("Seafarer").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *HealthRepo) Create(rec *model.HealthReexamination) error {
	return model.DB.Create(rec).Error
}

func (r *HealthRepo) GetByID(id int64) (*model.HealthReexamination, error) {
	var rec model.HealthReexamination
	if err := model.DB.Preload("Seafarer").First(&rec, id).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *HealthRepo) Update(rec *model.HealthReexamination) error {
	return model.DB.Save(rec).Error
}
