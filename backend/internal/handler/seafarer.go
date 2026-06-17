package handler

import (
	"strconv"
	"time"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/service"

	"github.com/gofiber/fiber/v2"
)

type SeafarerHandler struct {
	svc *service.SeafarerService
}

func NewSeafarerHandler(svc *service.SeafarerService) *SeafarerHandler {
	return &SeafarerHandler{svc: svc}
}

type seafarerDTO struct {
	Name     string `json:"name"`
	Gender   int8   `json:"gender"`
	Birthday string `json:"birthday"`
	IDNumber string `json:"id_number"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Rank     string `json:"rank"`
	Status   int8   `json:"status"`
}

func (d *seafarerDTO) toModel() model.Seafarer {
	s := model.Seafarer{
		Name:     d.Name,
		Gender:   d.Gender,
		IDNumber: d.IDNumber,
		Phone:    d.Phone,
		Email:    d.Email,
		Rank:     d.Rank,
		Status:   d.Status,
	}
	if d.Birthday != "" {
		if t, err := time.Parse("2006-01-02", d.Birthday); err == nil {
			s.Birthday = &t
		}
	}
	return s
}

func (h *SeafarerHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)
	keyword := c.Query("keyword")

	items, total, err := h.svc.List(page, pageSize, keyword)
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

func (h *SeafarerHandler) Get(c *fiber.Ctx) error {
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

func (h *SeafarerHandler) Create(c *fiber.Ctx) error {
	var body seafarerDTO
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	m := body.toModel()
	if err := h.svc.Create(&m); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(m))
}

func (h *SeafarerHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	var body seafarerDTO
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}
	m := body.toModel()
	m.ID = id

	if err := h.svc.Update(&m); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(m))
}

func (h *SeafarerHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	if err := h.svc.Delete(id); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(nil))
}
