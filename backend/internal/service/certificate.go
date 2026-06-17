package service

import (
	"errors"
	"time"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/repository"

	"github.com/sirupsen/logrus"
)

type CertificateService struct {
	typeRepo *repository.CertificateTypeRepo
	certRepo *repository.SeafarerCertRepo
	log      *logrus.Logger
}

func NewCertificateService(typeRepo *repository.CertificateTypeRepo, certRepo *repository.SeafarerCertRepo, log *logrus.Logger) *CertificateService {
	return &CertificateService{typeRepo: typeRepo, certRepo: certRepo, log: log}
}

func (s *CertificateService) ListTypes(page, pageSize int) ([]model.CertificateType, int64, error) {
	s.log.Debugf("CertificateService.ListTypes page=%d pageSize=%d", page, pageSize)
	return s.typeRepo.ListTypes(page, pageSize)
}

func (s *CertificateService) GetTypeByID(id int64) (*model.CertificateType, error) {
	return s.typeRepo.GetTypeByID(id)
}

func (s *CertificateService) CreateType(ct *model.CertificateType) error {
	s.log.Debugf("CertificateService.CreateType code=%s", ct.Code)
	if ct.Name == "" || ct.Code == "" {
		return errors.New("证书名称和编码不能为空")
	}
	return s.typeRepo.CreateType(ct)
}

func (s *CertificateService) UpdateType(ct *model.CertificateType) error {
	s.log.Debugf("CertificateService.UpdateType id=%d", ct.ID)
	existing, err := s.typeRepo.GetTypeByID(ct.ID)
	if err != nil {
		return err
	}
	existing.Name = ct.Name
	existing.Code = ct.Code
	existing.Description = ct.Description
	if ct.ValidityMonths != nil {
		existing.ValidityMonths = ct.ValidityMonths
	}
	existing.IsRequired = ct.IsRequired
	return s.typeRepo.UpdateType(existing)
}

func (s *CertificateService) DeleteType(id int64) error {
	return s.typeRepo.DeleteType(id)
}

func (s *CertificateService) ListBySeafarer(seafarerID int64) ([]model.SeafarerCertificate, error) {
	s.log.Debugf("CertificateService.ListBySeafarer seafarerID=%d", seafarerID)
	return s.certRepo.ListBySeafarerID(seafarerID)
}

func (s *CertificateService) GetCertByID(id int64) (*model.SeafarerCertificate, error) {
	return s.certRepo.GetCertByID(id)
}

func (s *CertificateService) CreateCert(sc *model.SeafarerCertificate) error {
	s.log.Debugf("CertificateService.CreateCert seafarerID=%d typeID=%d", sc.SeafarerID, sc.CertificateTypeID)
	if sc.SeafarerID <= 0 || sc.CertificateTypeID <= 0 {
		return errors.New("船员ID和证书类型ID不能为空")
	}
	if sc.CertNumber == "" {
		return errors.New("证书编号不能为空")
	}
	sc.Status = s.calcCertStatus(sc.ExpireDate)
	return s.certRepo.CreateCert(sc)
}

func (s *CertificateService) UpdateCert(sc *model.SeafarerCertificate) error {
	s.log.Debugf("CertificateService.UpdateCert id=%d", sc.ID)
	existing, err := s.certRepo.GetCertByID(sc.ID)
	if err != nil {
		return err
	}
	existing.SeafarerID = sc.SeafarerID
	existing.CertificateTypeID = sc.CertificateTypeID
	existing.CertNumber = sc.CertNumber
	existing.IssueDate = sc.IssueDate
	if sc.ExpireDate != nil {
		existing.ExpireDate = sc.ExpireDate
	}
	existing.CertImageURL = sc.CertImageURL
	existing.Status = s.calcCertStatus(existing.ExpireDate)
	return s.certRepo.UpdateCert(existing)
}

func (s *CertificateService) DeleteCert(id int64) error {
	return s.certRepo.DeleteCert(id)
}

func (s *CertificateService) calcCertStatus(expireDate *time.Time) int8 {
	if expireDate == nil {
		return model.CertStatusValid
	}
	now := time.Now()
	days := int(expireDate.Sub(now).Hours() / 24)
	if days <= 0 {
		return model.CertStatusExpired
	}
	if days <= 30 {
		return model.CertStatusExpiring
	}
	return model.CertStatusValid
}

func (s *CertificateService) RefreshAllCertStatus() (int, error) {
	s.log.Debug("CertificateService.RefreshAllCertStatus")
	certs90, err := s.certRepo.FindExpiring(90)
	if err != nil {
		return 0, err
	}
	updated := 0
	for _, c := range certs90 {
		newStatus := s.calcCertStatus(c.ExpireDate)
		if c.Status != newStatus {
			if err := s.certRepo.UpdateCertStatus(c.ID, newStatus); err != nil {
				s.log.Errorf("更新证书状态失败 certID=%d err=%v", c.ID, err)
				continue
			}
			updated++
		}
	}
	s.log.Debugf("证书状态刷新完成 更新%d条", updated)
	return updated, nil
}

func (s *CertificateService) FindExpiringCerts(days int) ([]model.SeafarerCertificate, error) {
	return s.certRepo.FindExpiring(days)
}
