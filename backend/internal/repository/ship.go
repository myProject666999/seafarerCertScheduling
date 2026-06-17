package repository

import (
	"seafarer-cert-scheduling/internal/model"
)

type ShipRepo struct{}

func (r *ShipRepo) List(page, pageSize int, keyword string) ([]model.Ship, int64, error) {
	var list []model.Ship
	var total int64
	db := model.DB.Model(&model.Ship{})
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("name LIKE ? OR imo_number LIKE ?", like, like)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes(paginate(page, pageSize)).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *ShipRepo) GetByID(id int64) (*model.Ship, error) {
	var s model.Ship
	if err := model.DB.Preload("Positions").First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ShipRepo) GetByIMONumber(imo string) (*model.Ship, error) {
	var s model.Ship
	if err := model.DB.Where("imo_number = ?", imo).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ShipRepo) Create(s *model.Ship) error {
	return model.DB.Create(s).Error
}

func (r *ShipRepo) Update(s *model.Ship) error {
	return model.DB.Save(s).Error
}

func (r *ShipRepo) Delete(id int64) error {
	return model.DB.Delete(&model.Ship{}, id).Error
}

type ShipPositionRepo struct{}

func (r *ShipPositionRepo) ListByShipID(shipID int64) ([]model.ShipPosition, error) {
	var list []model.ShipPosition
	if err := model.DB.Where("ship_id = ?", shipID).Preload("CertRequirements.CertificateType").Order("sort_order ASC, id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ShipPositionRepo) GetPositionByID(id int64) (*model.ShipPosition, error) {
	var sp model.ShipPosition
	if err := model.DB.Preload("CertRequirements.CertificateType").First(&sp, id).Error; err != nil {
		return nil, err
	}
	return &sp, nil
}

func (r *ShipPositionRepo) CreatePosition(sp *model.ShipPosition) error {
	return model.DB.Create(sp).Error
}

func (r *ShipPositionRepo) UpdatePosition(sp *model.ShipPosition) error {
	return model.DB.Save(sp).Error
}

func (r *ShipPositionRepo) DeletePosition(id int64) error {
	return model.DB.Delete(&model.ShipPosition{}, id).Error
}

func (r *ShipPositionRepo) GetOnboardCountForPosition(positionID int64) (int64, error) {
	var count int64
	if err := model.DB.Model(&model.SeafarerAssignment{}).Where("ship_position_id = ? AND status = ?", positionID, model.AssignmentStatusOnboard).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

type PositionCertReqRepo struct{}

func (r *PositionCertReqRepo) ListByPositionID(positionID int64) ([]model.PositionCertRequirement, error) {
	var list []model.PositionCertRequirement
	if err := model.DB.Where("ship_position_id = ?", positionID).Preload("CertificateType").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PositionCertReqRepo) CreateReq(req *model.PositionCertRequirement) error {
	return model.DB.Create(req).Error
}

func (r *PositionCertReqRepo) DeleteReq(id int64) error {
	return model.DB.Delete(&model.PositionCertRequirement{}, id).Error
}
