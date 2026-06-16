package handler

import (
	"strconv"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AlertHandler struct {
	svc *service.AlertService
}

func NewAlertHandler(svc *service.AlertService) *AlertHandler {
	return &AlertHandler{svc: svc}
}

func (h *AlertHandler) List(c *fiber.Ctx) error {
	seafarerID, _ := strconv.ParseInt(c.Query("seafarer_id"), 10, 64)
	level := int8(c.QueryInt("level"))
	isHandled := int8(c.QueryInt("is_handled"))

	items, total, err := h.svc.List(seafarerID, level, isHandled)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(model.PageResponse{
		Total: total,
		Items: items,
	}))
}

func (h *AlertHandler) HandleAlert(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	var body struct {
		Remark string `json:"remark"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	if err := h.svc.HandleAlert(id, body.Remark); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(nil))
}

func (h *AlertHandler) Stats(c *fiber.Ctx) error {
	stats := h.svc.GetAlertStats()
	return c.JSON(model.SuccessResponse(stats))
}

func (h *AlertHandler) RunScan(c *fiber.Ctx) error {
	count, err := h.svc.RunDailyAlertScan()
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(map[string]int{"count": count}))
}
