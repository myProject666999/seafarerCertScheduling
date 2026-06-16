package handler

import (
	"strconv"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AssignmentHandler struct {
	svc *service.AssignmentService
}

func NewAssignmentHandler(svc *service.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{svc: svc}
}

func (h *AssignmentHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)
	shipID, _ := strconv.ParseInt(c.Query("ship_id"), 10, 64)
	seafarerID, _ := strconv.ParseInt(c.Query("seafarer_id"), 10, 64)
	status := int8(c.QueryInt("status"))

	items, total, err := h.svc.List(page, pageSize, shipID, seafarerID, status)
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

func (h *AssignmentHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	item, err := h.svc.GetByID(id)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(item))
}

func (h *AssignmentHandler) Create(c *fiber.Ctx) error {
	var body model.SeafarerAssignment
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	if err := h.svc.Create(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(body))
}

func (h *AssignmentHandler) Disembark(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	var body struct {
		ActualDate string `json:"actual_date"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	if err := h.svc.Disembark(id, body.ActualDate); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(nil))
}
