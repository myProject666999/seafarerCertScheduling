package handler

import (
	"strconv"
	"time"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/service"

	"github.com/gofiber/fiber/v2"
)

type seafarerAssignmentDTO struct {
	SeafarerID            int64  `json:"seafarer_id"`
	ShipID                int64  `json:"ship_id"`
	ShipPositionID        int64  `json:"ship_position_id"`
	EmbarkDate            string `json:"embark_date"`
	ExpectedDisembarkDate string `json:"expected_disembark_date"`
	Status                int8   `json:"status"`
}

func (d *seafarerAssignmentDTO) toModel() (*model.SeafarerAssignment, error) {
	embarkDate, err := time.Parse("2006-01-02", d.EmbarkDate)
	if err != nil {
		return nil, err
	}

	assignment := &model.SeafarerAssignment{
		SeafarerID:     d.SeafarerID,
		ShipID:         d.ShipID,
		ShipPositionID: d.ShipPositionID,
		EmbarkDate:     embarkDate,
		Status:         d.Status,
	}

	if d.ExpectedDisembarkDate != "" {
		expectedDate, err := time.Parse("2006-01-02", d.ExpectedDisembarkDate)
		if err != nil {
			return nil, err
		}
		assignment.ExpectedDisembarkDate = &expectedDate
	}

	return assignment, nil
}

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
	var dto seafarerAssignmentDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	assignment, err := dto.toModel()
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "日期格式错误，请使用 YYYY-MM-DD 格式"))
	}

	if err := h.svc.Create(assignment); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(assignment))
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
