package router

import (
	"seafarer-cert-scheduling/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App,
	seafarerH *handler.SeafarerHandler,
	certH *handler.CertificateHandler,
	shipH *handler.ShipHandler,
	assignmentH *handler.AssignmentHandler,
	contractH *handler.ContractHandler,
	recordH *handler.RecordHandler,
	transferH *handler.TransferHandler,
	alertH *handler.AlertHandler,
) {
	api := app.Group("/api/v1")

	seafarers := api.Group("/seafarers")
	seafarers.Get("/", seafarerH.List)
	seafarers.Get("/:id", seafarerH.Get)
	seafarers.Post("/", seafarerH.Create)
	seafarers.Put("/:id", seafarerH.Update)
	seafarers.Delete("/:id", seafarerH.Delete)

	certTypes := api.Group("/certificate-types")
	certTypes.Get("/", certH.ListTypes)
	certTypes.Get("/:id", certH.GetType)
	certTypes.Post("/", certH.CreateType)
	certTypes.Put("/:id", certH.UpdateType)
	certTypes.Delete("/:id", certH.DeleteType)

	seafarerCerts := api.Group("/seafarer-certificates")
	seafarerCerts.Get("/", certH.ListBySeafarer)
	seafarerCerts.Get("/:id", certH.GetCert)
	seafarerCerts.Post("/", certH.CreateCert)
	seafarerCerts.Put("/:id", certH.UpdateCert)
	seafarerCerts.Delete("/:id", certH.DeleteCert)

	ships := api.Group("/ships")
	ships.Get("/", shipH.List)
	ships.Get("/:id", shipH.Get)
	ships.Post("/", shipH.Create)
	ships.Put("/:id", shipH.Update)
	ships.Delete("/:id", shipH.Delete)

	shipPositions := api.Group("/ship-positions")
	shipPositions.Get("/", shipH.ListPositions)
	shipPositions.Post("/", shipH.CreatePosition)
	shipPositions.Put("/:id", shipH.UpdatePosition)
	shipPositions.Delete("/:id", shipH.DeletePosition)

	positionCertReqs := api.Group("/position-cert-requirements")
	positionCertReqs.Get("/", shipH.ListCertReqs)
	positionCertReqs.Post("/", shipH.CreateCertReq)
	positionCertReqs.Delete("/:id", shipH.DeleteCertReq)

	assignments := api.Group("/assignments")
	assignments.Get("/", assignmentH.List)
	assignments.Get("/:id", assignmentH.Get)
	assignments.Post("/", assignmentH.Create)
	assignments.Post("/:id/disembark", assignmentH.Disembark)

	contracts := api.Group("/contracts")
	contracts.Get("/", contractH.List)
	contracts.Post("/", contractH.Create)
	contracts.Put("/:id", contractH.Update)

	embarkRecords := api.Group("/embark-records")
	embarkRecords.Get("/", recordH.ListEmbarkRecords)
	embarkRecords.Post("/", recordH.CreateEmbarkRecord)

	leaveRecords := api.Group("/leave-records")
	leaveRecords.Get("/", recordH.ListLeaveRecords)
	leaveRecords.Post("/", recordH.CreateLeaveRecord)
	leaveRecords.Post("/:id/end", recordH.EndLeave)

	healthRecords := api.Group("/health-records")
	healthRecords.Get("/", recordH.ListHealthRecords)
	healthRecords.Post("/", recordH.CreateHealthRecord)
	healthRecords.Put("/:id", recordH.UpdateHealthRecord)

	transfers := api.Group("/transfers")
	transfers.Get("/", transferH.List)
	transfers.Get("/:id", transferH.Get)
	transfers.Post("/", transferH.Create)
	transfers.Post("/:id/approve", transferH.Approve)
	transfers.Post("/:id/reject", transferH.Reject)
	transfers.Post("/:id/cancel", transferH.Cancel)

	alerts := api.Group("/alerts")
	alerts.Get("/", alertH.List)
	alerts.Get("/stats", alertH.Stats)
	alerts.Post("/:id/handle", alertH.HandleAlert)
	alerts.Post("/scan", alertH.RunScan)
}
