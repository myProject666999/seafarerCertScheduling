package handler

import (
	"strconv"
	"time"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/service"

	"github.com/gofiber/fiber/v2"
)

type voyageContractDTO struct {
	SeafarerID     int64  `json:"seafarer_id"`
	ShipID         int64  `json:"ship_id"`
	ContractNumber string `json:"contract_number"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	ActualEndDate  string `json:"actual_end_date"`
	Status         int8   `json:"status"`
	Remarks        string `json:"remarks"`
}

func (d *voyageContractDTO) toModel() (*model.VoyageContract, error) {
	startDate, err := time.Parse("2006-01-02", d.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", d.EndDate)
	if err != nil {
		return nil, err
	}

	contract := &model.VoyageContract{
		SeafarerID:     d.SeafarerID,
		ShipID:         d.ShipID,
		ContractNumber: d.ContractNumber,
		StartDate:      startDate,
		EndDate:        endDate,
		Status:         d.Status,
		Remarks:        d.Remarks,
	}

	if d.ActualEndDate != "" {
		actualEndDate, err := time.Parse("2006-01-02", d.ActualEndDate)
		if err != nil {
			return nil, err
		}
		contract.ActualEndDate = &actualEndDate
	}

	return contract, nil
}

type embarkDisembarkRecordDTO struct {
	SeafarerID int64  `json:"seafarer_id"`
	ShipID     int64  `json:"ship_id"`
	RecordType int8   `json:"record_type"`
	RecordDate string `json:"record_date"`
	Port       string `json:"port"`
	Reason     string `json:"reason"`
	Operator   string `json:"operator"`
}

func (d *embarkDisembarkRecordDTO) toModel() (*model.EmbarkDisembarkRecord, error) {
	recordDate, err := time.Parse("2006-01-02", d.RecordDate)
	if err != nil {
		return nil, err
	}

	return &model.EmbarkDisembarkRecord{
		SeafarerID: d.SeafarerID,
		ShipID:     d.ShipID,
		RecordType: d.RecordType,
		RecordDate: recordDate,
		Port:       d.Port,
		Reason:     d.Reason,
		Operator:   d.Operator,
	}, nil
}

type leaveRecordDTO struct {
	SeafarerID int64  `json:"seafarer_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	LeaveDays  *int   `json:"leave_days"`
	Status     int8   `json:"status"`
	Reason     string `json:"reason"`
}

func (d *leaveRecordDTO) toModel() (*model.LeaveRecord, error) {
	startDate, err := time.Parse("2006-01-02", d.StartDate)
	if err != nil {
		return nil, err
	}

	record := &model.LeaveRecord{
		SeafarerID: d.SeafarerID,
		StartDate:  startDate,
		LeaveDays:  d.LeaveDays,
		Status:     d.Status,
		Reason:     d.Reason,
	}

	if d.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", d.EndDate)
		if err != nil {
			return nil, err
		}
		record.EndDate = &endDate
	}

	return record, nil
}

type healthReexaminationDTO struct {
	SeafarerID      int64  `json:"seafarer_id"`
	ExamDate        string `json:"exam_date"`
	NextExamDate    string `json:"next_exam_date"`
	ExamResult      int8   `json:"exam_result"`
	ExamInstitution string `json:"exam_institution"`
	ReportURL       string `json:"report_url"`
	Restrictions    string `json:"restrictions"`
}

func (d *healthReexaminationDTO) toModel() (*model.HealthReexamination, error) {
	examDate, err := time.Parse("2006-01-02", d.ExamDate)
	if err != nil {
		return nil, err
	}

	record := &model.HealthReexamination{
		SeafarerID:      d.SeafarerID,
		ExamDate:        examDate,
		ExamResult:      d.ExamResult,
		ExamInstitution: d.ExamInstitution,
		ReportURL:       d.ReportURL,
		Restrictions:    d.Restrictions,
	}

	if d.NextExamDate != "" {
		nextExamDate, err := time.Parse("2006-01-02", d.NextExamDate)
		if err != nil {
			return nil, err
		}
		record.NextExamDate = &nextExamDate
	}

	return record, nil
}

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
	var dto voyageContractDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	contract, err := dto.toModel()
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "日期格式错误，请使用 YYYY-MM-DD 格式"))
	}

	if err := h.svc.Create(contract); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(contract))
}

func (h *ContractHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	var dto voyageContractDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	contract, err := dto.toModel()
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "日期格式错误，请使用 YYYY-MM-DD 格式"))
	}
	contract.ID = id

	if err := h.svc.Update(contract); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(contract))
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
	var dto embarkDisembarkRecordDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	record, err := dto.toModel()
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "日期格式错误，请使用 YYYY-MM-DD 格式"))
	}

	if err := h.svc.CreateEmbarkRecord(record); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(record))
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
	var dto leaveRecordDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	record, err := dto.toModel()
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "日期格式错误，请使用 YYYY-MM-DD 格式"))
	}

	if err := h.svc.CreateLeaveRecord(record); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(record))
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
	var dto healthReexaminationDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	record, err := dto.toModel()
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "日期格式错误，请使用 YYYY-MM-DD 格式"))
	}

	if err := h.svc.CreateHealthRecord(record); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(record))
}

func (h *RecordHandler) UpdateHealthRecord(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "无效的ID"))
	}

	var dto healthReexaminationDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.JSON(model.ErrorResponse(-1, "请求参数错误"))
	}

	record, err := dto.toModel()
	if err != nil {
		return c.JSON(model.ErrorResponse(-1, "日期格式错误，请使用 YYYY-MM-DD 格式"))
	}
	record.ID = id

	if err := h.svc.UpdateHealthRecord(record); err != nil {
		return c.JSON(model.ErrorResponse(-1, err.Error()))
	}
	return c.JSON(model.SuccessResponse(record))
}
