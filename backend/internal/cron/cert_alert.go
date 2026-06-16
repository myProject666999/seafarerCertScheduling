package cron

import (
	"seafarer-cert-scheduling/internal/service"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type CertAlertCron struct {
	cron     *cron.Cron
	alertSvc *service.AlertService
	log      *logrus.Logger
}

func NewCertAlertCron(alertSvc *service.AlertService, log *logrus.Logger) *CertAlertCron {
	return &CertAlertCron{
		cron:     cron.New(cron.WithSeconds()),
		alertSvc: alertSvc,
		log:      log,
	}
}

func (c *CertAlertCron) Start() error {
	_, err := c.cron.AddFunc("0 0 2 * * *", func() {
		c.log.Info("定时任务: 开始执行证书预警扫描")
		count, err := c.alertSvc.RunDailyAlertScan()
		if err != nil {
			c.log.Errorf("定时任务: 证书预警扫描失败: %v", err)
			return
		}
		c.log.Infof("定时任务: 证书预警扫描完成 新建%d条预警", count)
	})
	if err != nil {
		return err
	}
	c.cron.Start()
	c.log.Info("证书预警定时任务已启动 (每天凌晨2:00执行)")
	return nil
}

func (c *CertAlertCron) Stop() {
	c.cron.Stop()
	c.log.Info("证书预警定时任务已停止")
}
