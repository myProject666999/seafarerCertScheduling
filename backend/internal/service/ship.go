package service

import (
	"errors"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/repository"

	"github.com/sirupsen/logrus"
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

func (s *ShipService) CreateShip(ship *model.Ship) error {
	s.log.Debugf("ShipService.CreateShip name=%s", ship.Name)
	if ship.Name == "" {
		return errors.New("船名不能为空")
	}
	return s.shipRepo.Create(ship)
}

func (s *ShipService) UpdateShip(ship *model.Ship) error {
	return s.shipRepo.Update(ship)
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
	return s.positionRepo.UpdatePosition(sp)
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
