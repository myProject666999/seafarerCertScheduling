package service

import (
	"errors"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/repository"

	"github.com/sirupsen/logrus"
)

type AssignmentService struct {
	assignmentRepo *repository.AssignmentRepo
	seafarerRepo   *repository.SeafarerRepo
	log            *logrus.Logger
}

func NewAssignmentService(assignmentRepo *repository.AssignmentRepo, seafarerRepo *repository.SeafarerRepo, log *logrus.Logger) *AssignmentService {
	return &AssignmentService{assignmentRepo: assignmentRepo, seafarerRepo: seafarerRepo, log: log}
}

func (s *AssignmentService) List(page, pageSize int, shipID, seafarerID int64, status int8) ([]model.SeafarerAssignment, int64, error) {
	s.log.Debugf("AssignmentService.List page=%d shipID=%d seafarerID=%d status=%d", page, shipID, seafarerID, status)
	return s.assignmentRepo.List(page, pageSize, shipID, seafarerID, status)
}

func (s *AssignmentService) GetByID(id int64) (*model.SeafarerAssignment, error) {
	return s.assignmentRepo.GetByID(id)
}

func (s *AssignmentService) Create(a *model.SeafarerAssignment) error {
	s.log.Debugf("AssignmentService.Create seafarerID=%d shipID=%d positionID=%d", a.SeafarerID, a.ShipID, a.ShipPositionID)
	if a.SeafarerID <= 0 || a.ShipID <= 0 || a.ShipPositionID <= 0 {
		return errors.New("船员ID、船舶ID和岗位ID不能为空")
	}
	active, err := s.assignmentRepo.FindActiveBySeafarer(a.SeafarerID)
	if err == nil && active != nil {
		return errors.New("该船员已有在船分配记录，请先下船")
	}
	if err := s.assignmentRepo.Create(a); err != nil {
		return err
	}
	sf, err := s.seafarerRepo.GetByID(a.SeafarerID)
	if err != nil {
		return err
	}
	sf.Status = model.SeafarerStatusOnboard
	return s.seafarerRepo.Update(sf)
}

func (s *AssignmentService) Disembark(id int64, actualDate string) error {
	s.log.Debugf("AssignmentService.Disembark id=%d", id)
	a, err := s.assignmentRepo.GetByID(id)
	if err != nil {
		return err
	}
	if a.Status != model.AssignmentStatusOnboard {
		return errors.New("该分配记录已不在船")
	}
	now := parseDate(actualDate)
	a.ActualDisembarkDate = &now
	a.Status = model.AssignmentStatusDisembarked
	if err := s.assignmentRepo.Update(a); err != nil {
		return err
	}
	sf, err := s.seafarerRepo.GetByID(a.SeafarerID)
	if err != nil {
		return err
	}
	sf.Status = model.SeafarerStatusLeave
	return s.seafarerRepo.Update(sf)
}
