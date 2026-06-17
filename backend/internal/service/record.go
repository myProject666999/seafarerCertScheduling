package service

import (
	"time"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/repository"

	"github.com/sirupsen/logrus"
)

type ContractService struct {
	repo *repository.ContractRepo
	log  *logrus.Logger
}

func NewContractService(repo *repository.ContractRepo, log *logrus.Logger) *ContractService {
	return &ContractService{repo: repo, log: log}
}

func (s *ContractService) List(page, pageSize int, seafarerID, shipID int64) ([]model.VoyageContract, int64, error) {
	s.log.Debugf("ContractService.List page=%d seafarerID=%d shipID=%d", page, seafarerID, shipID)
	return s.repo.List(page, pageSize, seafarerID, shipID)
}

func (s *ContractService) Create(c *model.VoyageContract) error {
	s.log.Debugf("ContractService.Create number=%s", c.ContractNumber)
	return s.repo.Create(c)
}

func (s *ContractService) Update(c *model.VoyageContract) error {
	s.log.Debugf("ContractService.Update id=%d", c.ID)
	existing, err := s.repo.GetByID(c.ID)
	if err != nil {
		return err
	}
	existing.SeafarerID = c.SeafarerID
	existing.ShipID = c.ShipID
	existing.ContractNumber = c.ContractNumber
	existing.StartDate = c.StartDate
	existing.EndDate = c.EndDate
	if c.ActualEndDate != nil {
		existing.ActualEndDate = c.ActualEndDate
	}
	existing.Status = c.Status
	existing.Remarks = c.Remarks
	return s.repo.Update(existing)
}

type RecordService struct {
	embarkRepo *repository.EmbarkRecordRepo
	leaveRepo  *repository.LeaveRepo
	healthRepo *repository.HealthRepo
	log        *logrus.Logger
}

func NewRecordService(embarkRepo *repository.EmbarkRecordRepo, leaveRepo *repository.LeaveRepo, healthRepo *repository.HealthRepo, log *logrus.Logger) *RecordService {
	return &RecordService{embarkRepo: embarkRepo, leaveRepo: leaveRepo, healthRepo: healthRepo, log: log}
}

func (s *RecordService) ListEmbarkRecords(seafarerID int64) ([]model.EmbarkDisembarkRecord, error) {
	s.log.Debugf("RecordService.ListEmbarkRecords seafarerID=%d", seafarerID)
	return s.embarkRepo.List(seafarerID)
}

func (s *RecordService) CreateEmbarkRecord(r *model.EmbarkDisembarkRecord) error {
	s.log.Debugf("RecordService.CreateEmbarkRecord type=%d", r.RecordType)
	return s.embarkRepo.Create(r)
}

func (s *RecordService) ListLeaveRecords(seafarerID int64) ([]model.LeaveRecord, error) {
	s.log.Debugf("RecordService.ListLeaveRecords seafarerID=%d", seafarerID)
	return s.leaveRepo.List(seafarerID)
}

func (s *RecordService) CreateLeaveRecord(r *model.LeaveRecord) error {
	s.log.Debugf("RecordService.CreateLeaveRecord seafarerID=%d", r.SeafarerID)
	return s.leaveRepo.Create(r)
}

func (s *RecordService) EndLeave(id int64) error {
	s.log.Debugf("RecordService.EndLeave id=%d", id)
	lr, err := s.leaveRepo.GetActiveLeave(id)
	if err != nil {
		return err
	}
	now := time.Now()
	lr.EndDate = &now
	days := int(now.Sub(lr.StartDate).Hours() / 24)
	lr.LeaveDays = &days
	lr.Status = model.LeaveStatusEnded
	return s.leaveRepo.Update(lr)
}

func (s *RecordService) ListHealthRecords(seafarerID int64) ([]model.HealthReexamination, error) {
	s.log.Debugf("RecordService.ListHealthRecords seafarerID=%d", seafarerID)
	return s.healthRepo.List(seafarerID)
}

func (s *RecordService) CreateHealthRecord(r *model.HealthReexamination) error {
	s.log.Debugf("RecordService.CreateHealthRecord seafarerID=%d", r.SeafarerID)
	return s.healthRepo.Create(r)
}

func (s *RecordService) UpdateHealthRecord(r *model.HealthReexamination) error {
	s.log.Debugf("RecordService.UpdateHealthRecord id=%d", r.ID)
	existing, err := s.healthRepo.GetByID(r.ID)
	if err != nil {
		return err
	}
	existing.SeafarerID = r.SeafarerID
	existing.ExamDate = r.ExamDate
	if r.NextExamDate != nil {
		existing.NextExamDate = r.NextExamDate
	}
	existing.ExamResult = r.ExamResult
	existing.ExamInstitution = r.ExamInstitution
	existing.ReportURL = r.ReportURL
	existing.Restrictions = r.Restrictions
	return s.healthRepo.Update(existing)
}
