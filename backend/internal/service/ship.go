package service

import (
	"errors"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ShipService struct {
	shipRepo     *repository.ShipRepo
	positionRepo *repository.ShipPositionRepo
	reqRepo      *repository.PositionCertReqRepo
	log          *logrus.Logger
}

func NewShipService(shipRepo *repository.ShipRepo, positionRepo *repository.ShipPositionRepo, reqRepo *repository.PositionCertReqRepo, log *logrus.Logger) *ShipService {
	return &ShipService{shipRepo: shipRepo, positionRepo: positionRepo, reqRepo: reqRepo, log: log}
}

func (s *ShipService) ListShips(page, pageSize int, keyword string) ([]model.Ship, int64, error) {
	s.log.Debugf("ShipService.ListShips page=%d keyword=%s", page, keyword)
	return s.shipRepo.List(page, pageSize, keyword)
}

func (s *ShipService) GetShipByID(id int64) (*model.Ship, error) {
	return s.shipRepo.GetByID(id)
}

func normalizeEmptyStrToNil(s *string) *string {
	if s == nil {
		return nil
	}
	if *s == "" {
		return nil
	}
	return s
}

func (s *ShipService) CreateShip(ship *model.Ship) error {
	s.log.Debugf("ShipService.CreateShip name=%s", ship.Name)
	if ship.Name == "" {
		return errors.New("船名不能为空")
	}
	ship.IMONumber = normalizeEmptyStrToNil(ship.IMONumber)
	ship.MMSI = normalizeEmptyStrToNil(ship.MMSI)
	if ship.IMONumber != nil {
		existing, err := s.shipRepo.GetByIMONumber(*ship.IMONumber)
		if err == nil && existing != nil {
			return errors.New("该 IMO 号已存在")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	return s.shipRepo.Create(ship)
}

func (s *ShipService) UpdateShip(ship *model.Ship) error {
	s.log.Debugf("ShipService.UpdateShip id=%d", ship.ID)
	existing, err := s.shipRepo.GetByID(ship.ID)
	if err != nil {
		return err
	}
	existing.Name = ship.Name
	newIMO := normalizeEmptyStrToNil(ship.IMONumber)
	if newIMO != nil {
		found, err := s.shipRepo.GetByIMONumber(*newIMO)
		if err == nil && found != nil && found.ID != ship.ID {
			return errors.New("该 IMO 号已存在")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	existing.IMONumber = newIMO
	existing.MMSI = normalizeEmptyStrToNil(ship.MMSI)
	existing.ShipType = ship.ShipType
	if ship.GrossTonnage != nil {
		existing.GrossTonnage = ship.GrossTonnage
	}
	existing.FlagState = ship.FlagState
	existing.Status = ship.Status
	return s.shipRepo.Update(existing)
}

func (s *ShipService) DeleteShip(id int64) error {
	return s.shipRepo.Delete(id)
}

func (s *ShipService) ListPositions(shipID int64) ([]model.ShipPosition, error) {
	s.log.Debugf("ShipService.ListPositions shipID=%d", shipID)
	return s.positionRepo.ListByShipID(shipID)
}

func (s *ShipService) GetPositionByID(id int64) (*model.ShipPosition, error) {
	return s.positionRepo.GetPositionByID(id)
}

func (s *ShipService) CreatePosition(sp *model.ShipPosition) error {
	s.log.Debugf("ShipService.CreatePosition shipID=%d name=%s", sp.ShipID, sp.PositionName)
	if sp.ShipID <= 0 || sp.PositionName == "" {
		return errors.New("船舶ID和岗位名称不能为空")
	}
	return s.positionRepo.CreatePosition(sp)
}

func (s *ShipService) UpdatePosition(sp *model.ShipPosition) error {
	s.log.Debugf("ShipService.UpdatePosition id=%d", sp.ID)
	existing, err := s.positionRepo.GetPositionByID(sp.ID)
	if err != nil {
		return err
	}
	existing.ShipID = sp.ShipID
	existing.PositionName = sp.PositionName
	existing.Department = sp.Department
	existing.RequiredCount = sp.RequiredCount
	existing.SortOrder = sp.SortOrder
	return s.positionRepo.UpdatePosition(existing)
}

func (s *ShipService) DeletePosition(id int64) error {
	count, err := s.positionRepo.GetOnboardCountForPosition(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该岗位尚有在船人员，无法删除")
	}
	return s.positionRepo.DeletePosition(id)
}

func (s *ShipService) CreateCertReq(req *model.PositionCertRequirement) error {
	s.log.Debugf("ShipService.CreateCertReq positionID=%d typeID=%d", req.ShipPositionID, req.CertificateTypeID)
	return s.reqRepo.CreateReq(req)
}

func (s *ShipService) DeleteCertReq(id int64) error {
	return s.reqRepo.DeleteReq(id)
}

func (s *ShipService) ListCertReqs(positionID int64) ([]model.PositionCertRequirement, error) {
	return s.reqRepo.ListByPositionID(positionID)
}
