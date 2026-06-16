package repository

import (
	"seafarer-cert-scheduling/internal/model"

	"gorm.io/gorm"
)

type SeafarerRepo struct{}

func paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 20
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (r *SeafarerRepo) List(page, pageSize int, keyword string) ([]model.Seafarer, int64, error) {
	var list []model.Seafarer
	var total int64
	db := model.DB.Model(&model.Seafarer{})
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("name LIKE ? OR phone LIKE ? OR id_number LIKE ?", like, like, like)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes(paginate(page, pageSize)).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *SeafarerRepo) GetByID(id int64) (*model.Seafarer, error) {
	var s model.Seafarer
	if err := model.DB.First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SeafarerRepo) Create(s *model.Seafarer) error {
	return model.DB.Create(s).Error
}

func (r *SeafarerRepo) Update(s *model.Seafarer) error {
	return model.DB.Save(s).Error
}

func (r *SeafarerRepo) Delete(id int64) error {
	return model.DB.Delete(&model.Seafarer{}, id).Error
}
