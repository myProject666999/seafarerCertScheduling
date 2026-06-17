package handler

import (
	"strconv"
	"time"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/service"

	"github.com/gofiber/fiber/v2"
)

type seafarerCertificateDTO struct {
	SeafarerID        int64  `json:"seafarer_id"`
	CertificateTypeID int64  `json:"certificate_type_id"`
	CertNumber        string `json:"cert_number"`
	IssueDate         string `json:"issue_date"`
	ExpireDate        string `json:"expire_date"`
	CertImageURL      string `json:"cert_image_url"`
	Status            int8   `json:"status"`
}

func (d *seafarerCertificateDTO) toModel() (*model.SeafarerCertificate, error) {
	issueDate, err := time.Parse("2006-01-02", d.IssueDate)
	if err != nil {
		return nil, err
	}

	cert := &model.SeafarerCertificate{
		SeafarerID:        d.SeafarerID,
		CertificateTypeID: d.CertificateTypeID,
		CertNumber:        d.CertNumber,
		IssueDate:         issueDate,
		CertImageURL:      d.CertImageURL,
		Status:            d.Status,
	}

	if d.ExpireDate != "" {
		expireDate, err := time.Parse("2006-01-02", d.ExpireDate)
		if err != nil {
			return nil, err
		}
		cert.ExpireDate = &expireDate
	}

	return cert, nil
}

type CertificateHandler struct {
	svc *service.CertificateService
}

func NewCertificateHandler(svc *service.CertificateService) *CertificateHandler {
	return &CertificateHandler{svc: svc}
}

func (h *CertificateHandler) ListTypes(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 100)
	items, total, err := h.svc.ListTypes(page, pageSize)
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

func (h *CertificateHandler) GetType(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	item, err := h.svc.GetTypeByID(id)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(item))
}

func (h *CertificateHandler) CreateType(c *fiber.Ctx) error {
	var body model.CertificateType
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	if err := h.svc.CreateType(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(body))
}

func (h *CertificateHandler) UpdateType(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	var body model.CertificateType
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}
	body.ID = id

	if err := h.svc.UpdateType(&body); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(body))
}

func (h *CertificateHandler) DeleteType(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	if err := h.svc.DeleteType(id); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(nil))
}

func (h *CertificateHandler) ListBySeafarer(c *fiber.Ctx) error {
	seafarerID, err := strconv.ParseInt(c.Params("seafarer_id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的船员ID"))
	}

	items, err := h.svc.ListBySeafarer(seafarerID)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(items))
}

func (h *CertificateHandler) GetCert(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	item, err := h.svc.GetCertByID(id)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(item))
}

func (h *CertificateHandler) CreateCert(c *fiber.Ctx) error {
	var dto seafarerCertificateDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	cert, err := dto.toModel()
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "日期格式错误，请使用 YYYY-MM-DD 格式"))
	}

	if err := h.svc.CreateCert(cert); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(cert))
}

func (h *CertificateHandler) UpdateCert(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	var dto seafarerCertificateDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	cert, err := dto.toModel()
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "日期格式错误，请使用 YYYY-MM-DD 格式"))
	}
	cert.ID = id

	if err := h.svc.UpdateCert(cert); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(cert))
}

func (h *CertificateHandler) DeleteCert(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	if err := h.svc.DeleteCert(id); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(nil))
}
