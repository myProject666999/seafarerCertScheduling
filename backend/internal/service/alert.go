package service

import (
	"fmt"
	"time"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/repository"

	"github.com/sirupsen/logrus"
)

type AlertService struct {
	alertRepo *repository.AlertRepo
	certRepo  *repository.SeafarerCertRepo
	certSvc   *CertificateService
	log       *logrus.Logger
}

func NewAlertService(alertRepo *repository.AlertRepo, certRepo *repository.SeafarerCertRepo, certSvc *CertificateService, log *logrus.Logger) *AlertService {
	return &AlertService{alertRepo: alertRepo, certRepo: certRepo, certSvc: certSvc, log: log}
}

func (s *AlertService) List(seafarerID int64, level, isHandled int8) ([]model.CertAlert, int64, error) {
	s.log.Debugf("AlertService.List seafarerID=%d level=%d isHandled=%d", seafarerID, level, isHandled)
	return s.alertRepo.List(seafarerID, level, isHandled)
}

func (s *AlertService) HandleAlert(id int64, remark string) error {
	s.log.Debugf("AlertService.HandleAlert id=%d", id)
	var alert model.CertAlert
	alert.ID = id
	alert.IsHandled = 1
	alert.HandleRemark = remark
	return s.alertRepo.Update(&alert)
}

func (s *AlertService) RunDailyAlertScan() (int, error) {
	s.log.Debug("AlertService.RunDailyAlertScan 开始执行证书预警扫描")

	_, _ = s.certSvc.RefreshAllCertStatus()

	totalCreated := 0

	levels := []struct {
		days  int
		level int8
	}{
		{90, model.AlertLevel90Days},
		{60, model.AlertLevel60Days},
		{30, model.AlertLevel30Days},
	}

	for _, l := range levels {
		certs, err := s.certSvc.FindExpiringCerts(l.days)
		if err != nil {
			s.log.Errorf("查询%d天内到期证书失败: %v", l.days, err)
			continue
		}

		for _, cert := range certs {
			if cert.ExpireDate == nil {
				continue
			}

			exists, err := s.alertRepo.ExistsToday(cert.ID, l.level)
			if err != nil {
				s.log.Errorf("检查预警是否存在失败 certID=%d: %v", cert.ID, err)
				continue
			}
			if exists {
				continue
			}

			daysRemaining := int(cert.ExpireDate.Sub(time.Now()).Hours() / 24)
			if daysRemaining < 0 {
				daysRemaining = 0
			}

			higherAlertExists := false
			if l.level == model.AlertLevel90Days {
				e60, _ := s.alertRepo.ExistsToday(cert.ID, model.AlertLevel60Days)
				e30, _ := s.alertRepo.ExistsToday(cert.ID, model.AlertLevel30Days)
				if e60 || e30 {
					higherAlertExists = true
				}
			} else if l.level == model.AlertLevel60Days {
				e30, _ := s.alertRepo.ExistsToday(cert.ID, model.AlertLevel30Days)
				if e30 {
					higherAlertExists = true
				}
			}

			if higherAlertExists {
				continue
			}

			alert := &model.CertAlert{
				SeafarerCertificateID: cert.ID,
				SeafarerID:            cert.SeafarerID,
				AlertLevel:            l.level,
				AlertDate:             time.Now(),
				ExpireDate:            *cert.ExpireDate,
				DaysRemaining:         daysRemaining,
				IsHandled:             0,
			}
			if err := s.alertRepo.Create(alert); err != nil {
				s.log.Errorf("创建预警记录失败 certID=%d level=%d: %v", cert.ID, l.level, err)
				continue
			}
			totalCreated++
			s.log.Debugf("创建预警: 船员证书ID=%d 级别=%d 剩余%d天", cert.ID, l.level, daysRemaining)
		}
	}

	s.log.Debugf("证书预警扫描完成 新建预警%d条", totalCreated)
	return totalCreated, nil
}

func (s *AlertService) GetAlertStats() map[string]int64 {
	s.log.Debug("AlertService.GetAlertStats")
	stats := make(map[string]int64)

	for _, level := range []int8{model.AlertLevel90Days, model.AlertLevel60Days, model.AlertLevel30Days} {
		_, total, err := s.alertRepo.List(0, level, 0)
		if err != nil {
			continue
		}
		key := fmt.Sprintf("level_%d", level)
		stats[key] = total
	}

	unhandled, _, _ := s.alertRepo.List(0, 0, 0)
	stats["total_unhandled"] = int64(len(unhandled))

	return stats
}
