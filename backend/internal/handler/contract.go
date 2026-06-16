package handler

import (
	"strconv"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ContractHandler struct {
	svc *service.ContractService
}

func NewContractHandler(svc *service.ContractService) *ContractHandler {
	return &ContractHandler{svc: svc}
}

func (h *ContractHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)
	seafarerID, _ := strconv.ParseInt(c.Query("seafarer_id"), 10, 64)
	shipID, _ := strconv.ParseInt(c.Query("ship_id"), 10, 64)

	items, total, err := h.svc.List(page, pageSize, seafarerID, shipID)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(model.PageResponse{
		Total: total,
		Page:  page,
		Size:  pageSize,
		Items: items,
	}))
}

func (h *ContractHandler) Create(c *fiber.Ctx) error {
	var body model.VoyageContract
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	if err := h.svc.Create(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(body))
}

func (h *ContractHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	var body model.VoyageContract
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}
	body.ID = id

	if err := h.svc.Update(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(body))
}

type RecordHandler struct {
	svc *service.RecordService
}

func NewRecordHandler(svc *service.RecordService) *RecordHandler {
	return &RecordHandler{svc: svc}
}

func (h *RecordHandler) ListEmbarkRecords(c *fiber.Ctx) error {
	seafarerID, _ := strconv.ParseInt(c.Query("seafarer_id"), 10, 64)

	items, err := h.svc.ListEmbarkRecords(seafarerID)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(items))
}

func (h *RecordHandler) CreateEmbarkRecord(c *fiber.Ctx) error {
	var body model.EmbarkDisembarkRecord
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	if err := h.svc.CreateEmbarkRecord(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(body))
}

func (h *RecordHandler) ListLeaveRecords(c *fiber.Ctx) error {
	seafarerID, _ := strconv.ParseInt(c.Query("seafarer_id"), 10, 64)

	items, err := h.svc.ListLeaveRecords(seafarerID)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(items))
}

func (h *RecordHandler) CreateLeaveRecord(c *fiber.Ctx) error {
	var body model.LeaveRecord
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	if err := h.svc.CreateLeaveRecord(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(body))
}

func (h *RecordHandler) EndLeave(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	if err := h.svc.EndLeave(id); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(nil))
}

func (h *RecordHandler) ListHealthRecords(c *fiber.Ctx) error {
	seafarerID, _ := strconv.ParseInt(c.Query("seafarer_id"), 10, 64)

	items, err := h.svc.ListHealthRecords(seafarerID)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(items))
}

func (h *RecordHandler) CreateHealthRecord(c *fiber.Ctx) error {
	var body model.HealthReexamination
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	if err := h.svc.CreateHealthRecord(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(body))
}

func (h *RecordHandler) UpdateHealthRecord(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	var body model.HealthReexamination
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}
	body.ID = id

	if err := h.svc.UpdateHealthRecord(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(body))
}
