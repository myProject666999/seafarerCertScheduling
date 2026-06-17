package repository

import (
	"time"

	"seafarer-cert-scheduling/internal/model"
)

type TransferRepo struct{}

func (r *TransferRepo) List(page, pageSize int, status int8) ([]model.TransferRequest, int64, error) {
	var list []model.TransferRequest
	var total int64
	db := model.DB.Model(&model.TransferRequest{})
	if status >= 0 {
		db = db.Where("status = ?", status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes(paginate(page, pageSize)).
		Preload("Seafarer").Preload("FromShip").Preload("ToShip").
		Preload("FromPosition").Preload("ToPosition").Preload("ReplacementSeafarer").
		Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *TransferRepo) GetByID(id int64) (*model.TransferRequest, error) {
	var t model.TransferRequest
	if err := model.DB.
		Preload("Seafarer").Preload("FromShip").Preload("ToShip").
		Preload("FromPosition").Preload("ToPosition").Preload("ReplacementSeafarer").
		First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TransferRepo) Create(t *model.TransferRequest) error {
	return model.DB.Create(t).Error
}

func (r *TransferRepo) Update(t *model.TransferRequest) error {
	return model.DB.Save(t).Error
}

type AlertRepo struct{}

func (r *AlertRepo) GetByID(id int64) (*model.CertAlert, error) {
	var a model.CertAlert
	if err := model.DB.Preload("SeafarerCertificate.CertificateType").Preload("Seafarer").First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AlertRepo) List(seafarerID int64, level, isHandled int8) ([]model.CertAlert, int64, error) {
	var list []model.CertAlert
	var total int64
	db := model.DB.Model(&model.CertAlert{})
	if seafarerID > 0 {
		db = db.Where("seafarer_id = ?", seafarerID)
	}
	if level > 0 {
		db = db.Where("alert_level = ?", level)
	}
	if isHandled >= 0 {
		db = db.Where("is_handled = ?", isHandled)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Preload("SeafarerCertificate.CertificateType").Preload("Seafarer").Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *AlertRepo) Create(a *model.CertAlert) error {
	return model.DB.Create(a).Error
}

func (r *AlertRepo) Update(a *model.CertAlert) error {
	return model.DB.Save(a).Error
}

func (r *AlertRepo) ExistsToday(certID int64, level int8) (bool, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	if err := model.DB.Model(&model.CertAlert{}).Where("seafarer_certificate_id = ? AND alert_level = ? AND alert_date >= ? AND alert_date < ?", certID, level, today, tomorrow).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
