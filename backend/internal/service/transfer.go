package service

import (
	"errors"
	"fmt"
	"time"

	"seafarer-cert-scheduling/internal/model"
	"seafarer-cert-scheduling/internal/repository"

	"github.com/sirupsen/logrus"
)

type TransferService struct {
	transferRepo  *repository.TransferRepo
	assignmentRepo *repository.AssignmentRepo
	positionRepo  *repository.ShipPositionRepo
	certRepo      *repository.SeafarerCertRepo
	seafarerRepo  *repository.SeafarerRepo
	reqRepo       *repository.PositionCertReqRepo
	log           *logrus.Logger
}

func NewTransferService(
	transferRepo *repository.TransferRepo,
	assignmentRepo *repository.AssignmentRepo,
	positionRepo *repository.ShipPositionRepo,
	certRepo *repository.SeafarerCertRepo,
	seafarerRepo *repository.SeafarerRepo,
	reqRepo *repository.PositionCertReqRepo,
	log *logrus.Logger,
) *TransferService {
	return &TransferService{
		transferRepo:  transferRepo,
		assignmentRepo: assignmentRepo,
		positionRepo:  positionRepo,
		certRepo:      certRepo,
		seafarerRepo:  seafarerRepo,
		reqRepo:       reqRepo,
		log:           log,
	}
}

type ValidationResult struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors"`
}

func (s *TransferService) List(page, pageSize int, status int8) ([]model.TransferRequest, int64, error) {
	s.log.Debugf("TransferService.List page=%d status=%d", page, status)
	return s.transferRepo.List(page, pageSize, status)
}

func (s *TransferService) GetByID(id int64) (*model.TransferRequest, error) {
	return s.transferRepo.GetByID(id)
}

func (s *TransferService) Create(t *model.TransferRequest) error {
	s.log.Debugf("TransferService.Create seafarerID=%d fromShip=%d toShip=%d", t.SeafarerID, t.FromShipID, t.ToShipID)
	if t.SeafarerID <= 0 || t.FromShipID <= 0 || t.ToShipID <= 0 || t.FromPositionID <= 0 || t.ToPositionID <= 0 {
		return errors.New("船员ID、原船ID、目标船ID、原岗位ID、目标岗位ID不能为空")
	}
	if t.FromShipID == t.ToShipID {
		return errors.New("原船和目标船不能相同")
	}
	if t.Reason == "" {
		return errors.New("调动原因不能为空")
	}

	activeAssign, err := s.assignmentRepo.FindActiveBySeafarer(t.SeafarerID)
	if err != nil || activeAssign == nil {
		return errors.New("该船员当前不在船上，无法发起调动")
	}
	if activeAssign.ShipID != t.FromShipID {
		return errors.New("船员当前所在船舶与原船ID不匹配")
	}
	if activeAssign.ShipPositionID != t.FromPositionID {
		return errors.New("船员当前岗位与原岗位ID不匹配")
	}

	fromResult := s.validateFromShip(t)
	toResult := s.validateToShip(t)

	validFrom := int8(0)
	validTo := int8(0)
	if fromResult.IsValid {
		validFrom = 1
	}
	if toResult.IsValid {
		validTo = 1
	}
	t.FromShipValid = &validFrom
	t.ToShipValid = &validTo

	if !fromResult.IsValid || !toResult.IsValid {
		allErrors := append(fromResult.Errors, toResult.Errors...)
		s.log.Debugf("调动校验不通过: %v", allErrors)
	}

	return s.transferRepo.Create(t)
}

func (s *TransferService) Approve(id int64, approver, remark string) error {
	s.log.Debugf("TransferService.Approve id=%d approver=%s", id, approver)
	t, err := s.transferRepo.GetByID(id)
	if err != nil {
		return err
	}
	if t.Status != model.TransferStatusPending {
		return errors.New("只能审批待审批状态的调动申请")
	}

	fromResult := s.validateFromShip(t)
	toResult := s.validateToShip(t)

	validFrom := int8(0)
	validTo := int8(0)
	if fromResult.IsValid {
		validFrom = 1
	}
	if toResult.IsValid {
		validTo = 1
	}
	t.FromShipValid = &validFrom
	t.ToShipValid = &validTo

	if !fromResult.IsValid || !toResult.IsValid {
		allErrors := append(fromResult.Errors, toResult.Errors...)
		t.Status = model.TransferStatusRejected
		t.Approver = approver
		t.ApproveRemark = fmt.Sprintf("校验不通过: %v", allErrors)
		now := time.Now()
		t.ApprovedAt = &now
		_ = s.transferRepo.Update(t)
		return fmt.Errorf("审批不通过: %v", allErrors)
	}

	if err := s.executeTransfer(t); err != nil {
		return fmt.Errorf("执行调动失败: %v", err)
	}

	now := time.Now()
	t.Status = model.TransferStatusApproved
	t.Approver = approver
	t.ApproveRemark = remark
	t.ApprovedAt = &now
	return s.transferRepo.Update(t)
}

func (s *TransferService) Reject(id int64, approver, remark string) error {
	s.log.Debugf("TransferService.Reject id=%d", id)
	t, err := s.transferRepo.GetByID(id)
	if err != nil {
		return err
	}
	if t.Status != model.TransferStatusPending {
		return errors.New("只能拒绝待审批状态的调动申请")
	}
	now := time.Now()
	t.Status = model.TransferStatusRejected
	t.Approver = approver
	t.ApproveRemark = remark
	t.ApprovedAt = &now
	return s.transferRepo.Update(t)
}

func (s *TransferService) Cancel(id int64) error {
	t, err := s.transferRepo.GetByID(id)
	if err != nil {
		return err
	}
	if t.Status != model.TransferStatusPending {
		return errors.New("只能取消待审批状态的调动申请")
	}
	t.Status = model.TransferStatusCancelled
	return s.transferRepo.Update(t)
}

func (s *TransferService) validateFromShip(t *model.TransferRequest) ValidationResult {
	result := ValidationResult{IsValid: true}

	fromPosition, err := s.positionRepo.GetPositionByID(t.FromPositionID)
	if err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, "原岗位不存在")
		return result
	}

	currentOnboard, err := s.positionRepo.GetOnboardCountForPosition(t.FromPositionID)
	if err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, "查询原岗位在船人数失败")
		return result
	}

	afterRemove := currentOnboard - 1

	if t.ReplacementSeafarerID != nil && *t.ReplacementSeafarerID > 0 {
		repActive, err := s.assignmentRepo.FindActiveBySeafarer(*t.ReplacementSeafarerID)
		if err != nil || repActive == nil {
			result.IsValid = false
			result.Errors = append(result.Errors, "替换船员当前不在可分配状态")
		} else if repActive.ShipID == t.FromShipID {
			result.IsValid = false
			result.Errors = append(result.Errors, "替换船员已在原船上，不能作为补位")
		}

		if repActive != nil {
			repCerts, err := s.certRepo.ListBySeafarerID(*t.ReplacementSeafarerID)
			if err != nil {
				result.IsValid = false
				result.Errors = append(result.Errors, "查询替换船员证书失败")
			} else {
				certReqs, _ := s.reqRepo.ListByPositionID(t.FromPositionID)
				for _, req := range certReqs {
					if req.IsMandatory == 1 {
						found := false
						for _, c := range repCerts {
							if c.CertificateTypeID == req.CertificateTypeID && c.Status != model.CertStatusExpired {
								found = true
								break
							}
						}
						if !found {
							result.IsValid = false
							result.Errors = append(result.Errors, fmt.Sprintf("替换船员缺少必要证书(要求ID=%d)", req.CertificateTypeID))
						}
					}
				}
			}
		}
		afterRemove = currentOnboard
	}

	if afterRemove < int64(fromPosition.RequiredCount) {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("原船岗位编制不足: 需要%d人，调动后仅剩%d人", fromPosition.RequiredCount, afterRemove))
	}

	return result
}

func (s *TransferService) validateToShip(t *model.TransferRequest) ValidationResult {
	result := ValidationResult{IsValid: true}

	toPosition, err := s.positionRepo.GetPositionByID(t.ToPositionID)
	if err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, "目标岗位不存在")
		return result
	}

	currentOnboard, err := s.positionRepo.GetOnboardCountForPosition(t.ToPositionID)
	if err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, "查询目标岗位在船人数失败")
		return result
	}

	if currentOnboard >= int64(toPosition.RequiredCount) {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("目标岗位已满编: 编制%d人，当前已有%d人", toPosition.RequiredCount, currentOnboard))
	}

	seafarerCerts, err := s.certRepo.ListBySeafarerID(t.SeafarerID)
	if err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, "查询调动船员证书失败")
		return result
	}

	certReqs, _ := s.reqRepo.ListByPositionID(t.ToPositionID)
	for _, req := range certReqs {
		if req.IsMandatory == 1 {
			found := false
			for _, c := range seafarerCerts {
				if c.CertificateTypeID == req.CertificateTypeID && c.Status != model.CertStatusExpired {
					found = true
					break
				}
			}
			if !found {
				result.IsValid = false
				result.Errors = append(result.Errors, fmt.Sprintf("调动船员缺少目标岗位必要证书(要求ID=%d)", req.CertificateTypeID))
			}
		}
	}

	return result
}

func (s *TransferService) executeTransfer(t *model.TransferRequest) error {
	activeAssign, err := s.assignmentRepo.FindActiveBySeafarer(t.SeafarerID)
	if err != nil {
		return err
	}

	now := time.Now()
	activeAssign.Status = model.AssignmentStatusDisembarked
	activeAssign.ActualDisembarkDate = &now
	if err := s.assignmentRepo.Update(activeAssign); err != nil {
		return err
	}

	newAssign := &model.SeafarerAssignment{
		SeafarerID:     t.SeafarerID,
		ShipID:         t.ToShipID,
		ShipPositionID: t.ToPositionID,
		EmbarkDate:     now,
		Status:         model.AssignmentStatusOnboard,
	}
	if err := s.assignmentRepo.Create(newAssign); err != nil {
		return err
	}

	if t.ReplacementSeafarerID != nil && *t.ReplacementSeafarerID > 0 {
		repActive, err := s.assignmentRepo.FindActiveBySeafarer(*t.ReplacementSeafarerID)
		if err == nil && repActive != nil {
			repActive.Status = model.AssignmentStatusDisembarked
			repActive.ActualDisembarkDate = &now
			_ = s.assignmentRepo.Update(repActive)
		}

		repNewAssign := &model.SeafarerAssignment{
			SeafarerID:     *t.ReplacementSeafarerID,
			ShipID:         t.FromShipID,
			ShipPositionID: t.FromPositionID,
			EmbarkDate:     now,
			Status:         model.AssignmentStatusOnboard,
		}
		if err := s.assignmentRepo.Create(repNewAssign); err != nil {
			s.log.Errorf("创建替换船员分配记录失败: %v", err)
		}
	}

	s.log.Debugf("调动执行成功: 船员%d 从船%d(%d) 调至 船%d(%d)", t.SeafarerID, t.FromShipID, t.FromPositionID, t.ToShipID, t.ToPositionID)
	return nil
}
