package repository

import (
	"time"

	"seafarer-cert-scheduling/internal/model"
)

type CertificateTypeRepo struct{}

func (r *CertificateTypeRepo) ListTypes() ([]model.CertificateType, error) {
	var list []model.CertificateType
	if err := model.DB.Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *CertificateTypeRepo) GetTypeByID(id int64) (*model.CertificateType, error) {
	var ct model.CertificateType
	if err := model.DB.First(&ct, id).Error; err != nil {
		return nil, err
	}
	return &ct, nil
}

func (r *CertificateTypeRepo) CreateType(ct *model.CertificateType) error {
	return model.DB.Create(ct).Error
}

func (r *CertificateTypeRepo) UpdateType(ct *model.CertificateType) error {
	return model.DB.Save(ct).Error
}

func (r *CertificateTypeRepo) DeleteType(id int64) error {
	return model.DB.Delete(&model.CertificateType{}, id).Error
}

type SeafarerCertRepo struct{}

func (r *SeafarerCertRepo) ListBySeafarerID(seafarerID int64) ([]model.SeafarerCertificate, error) {
	var list []model.SeafarerCertificate
	if err := model.DB.Where("seafarer_id = ?", seafarerID).Preload("CertificateType").Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *SeafarerCertRepo) GetCertByID(id int64) (*model.SeafarerCertificate, error) {
	var sc model.SeafarerCertificate
	if err := model.DB.Preload("CertificateType").First(&sc, id).Error; err != nil {
		return nil, err
	}
	return &sc, nil
}

func (r *SeafarerCertRepo) CreateCert(sc *model.SeafarerCertificate) error {
	return model.DB.Create(sc).Error
}

func (r *SeafarerCertRepo) UpdateCert(sc *model.SeafarerCertificate) error {
	return model.DB.Save(sc).Error
}

func (r *SeafarerCertRepo) DeleteCert(id int64) error {
	return model.DB.Delete(&model.SeafarerCertificate{}, id).Error
}

func (r *SeafarerCertRepo) FindExpiring(days int) ([]model.SeafarerCertificate, error) {
	var list []model.SeafarerCertificate
	now := time.Now()
	deadline := now.AddDate(0, 0, days)
	if err := model.DB.Where("expire_date <= ? AND expire_date > ? AND status != ?", deadline, now, model.CertStatusExpired).
		Preload("CertificateType").Preload("Seafarer").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *SeafarerCertRepo) UpdateCertStatus(id int64, status int8) error {
	return model.DB.Model(&model.SeafarerCertificate{}).Where("id = ?", id).Update("status", status).Error
}
