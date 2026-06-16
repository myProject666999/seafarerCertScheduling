package main

import (
	"os"
	"os/signal"
	"syscall"

	"seafarer-cert-scheduling/internal/config"
	"seafarer-cert-scheduling/internal/cron"
	"seafarer-cert-scheduling/internal/handler"
	"seafarer-cert-scheduling/internal/middleware"
	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/repository"
	"seafarer-cert-scheduling/internal/router"
	"seafarer-cert-scheduling/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()
	log := config.InitLogger()

	log.Debug("正在初始化数据库连接...")
	if err := model.InitDB(cfg.Database.DSN(), log); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	log.Debug("数据库连接成功")

	seafarerRepo := &repository.SeafarerRepo{}
	certTypeRepo := &repository.CertificateTypeRepo{}
	seafarerCertRepo := &repository.SeafarerCertRepo{}
	shipRepo := &repository.ShipRepo{}
	positionRepo := &repository.ShipPositionRepo{}
	positionCertReqRepo := &repository.PositionCertReqRepo{}
	assignmentRepo := &repository.AssignmentRepo{}
	contractRepo := &repository.ContractRepo{}
	embarkRepo := &repository.EmbarkRecordRepo{}
	leaveRepo := &repository.LeaveRepo{}
	healthRepo := &repository.HealthRepo{}
	transferRepo := &repository.TransferRepo{}
	alertRepo := &repository.AlertRepo{}

	seafarerSvc := service.NewSeafarerService(seafarerRepo, log)
	certSvc := service.NewCertificateService(certTypeRepo, seafarerCertRepo, log)
	shipSvc := service.NewShipService(shipRepo, positionRepo, positionCertReqRepo, log)
	assignmentSvc := service.NewAssignmentService(assignmentRepo, seafarerRepo, log)
	contractSvc := service.NewContractService(contractRepo, log)
	recordSvc := service.NewRecordService(embarkRepo, leaveRepo, healthRepo, log)
	transferSvc := service.NewTransferService(transferRepo, assignmentRepo, positionRepo, seafarerCertRepo, seafarerRepo, positionCertReqRepo, log)
	alertSvc := service.NewAlertService(alertRepo, seafarerCertRepo, certSvc, log)

	seafarerH := handler.NewSeafarerHandler(seafarerSvc)
	certH := handler.NewCertificateHandler(certSvc)
	shipH := handler.NewShipHandler(shipSvc)
	assignmentH := handler.NewAssignmentHandler(assignmentSvc)
	contractH := handler.NewContractHandler(contractSvc)
	recordH := handler.NewRecordHandler(recordSvc)
	transferH := handler.NewTransferHandler(transferSvc)
	alertH := handler.NewAlertHandler(alertSvc)

	certAlertCron := cron.NewCertAlertCron(alertSvc, log)
	if err := certAlertCron.Start(); err != nil {
		log.Errorf("启动证书预警定时任务失败: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: "SeafarerCertScheduling",
	})

	app.Use(middleware.RequestLogger(log))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(model.SuccessResponse("ok"))
	})

	router.Setup(app, seafarerH, certH, shipH, assignmentH, contractH, recordH, transferH, alertH)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Info("正在关闭服务...")
		certAlertCron.Stop()
		_ = app.Shutdown()
	}()

	log.Infof("服务启动于 http://localhost%s", cfg.Server.Port)
	if err := app.Listen(cfg.Server.Port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
