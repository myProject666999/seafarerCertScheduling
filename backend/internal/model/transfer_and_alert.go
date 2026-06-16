package model

import "time"

type TransferRequest struct {
	ID                  int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	SeafarerID          int64      `json:"seafarer_id" gorm:"not null;index"`
	FromShipID          int64      `json:"from_ship_id" gorm:"not null;index"`
	ToShipID            int64      `json:"to_ship_id" gorm:"not null;index"`
	FromPositionID      int64      `json:"from_position_id" gorm:"not null"`
	ToPositionID        int64      `json:"to_position_id" gorm:"not null"`
	ReplacementSeafarerID *int64   `json:"replacement_seafarer_id"`
	Reason              string     `json:"reason" gorm:"size:500;not null"`
	Status              int8       `json:"status" gorm:"not null;default:0;index"`
	Approver            string     `json:"approver" gorm:"size:50"`
	ApproveRemark       string     `json:"approve_remark" gorm:"size:500"`
	ApprovedAt          *time.Time `json:"approved_at"`
	FromShipValid       *int8      `json:"from_ship_valid"`
	ToShipValid         *int8      `json:"to_ship_valid"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`

	Seafarer           *Seafarer     `json:"seafarer,omitempty" gorm:"foreignKey:SeafarerID"`
	FromShip           *Ship         `json:"from_ship,omitempty" gorm:"foreignKey:FromShipID"`
	ToShip             *Ship         `json:"to_ship,omitempty" gorm:"foreignKey:ToShipID"`
	FromPosition       *ShipPosition `json:"from_position,omitempty" gorm:"foreignKey:FromPositionID"`
	ToPosition         *ShipPosition `json:"to_position,omitempty" gorm:"foreignKey:ToPositionID"`
	ReplacementSeafarer *Seafarer    `json:"replacement_seafarer,omitempty" gorm:"foreignKey:ReplacementSeafarerID"`
}

func (TransferRequest) TableName() string {
	return "transfer_request"
}

const (
	TransferStatusPending = iota
	TransferStatusApproved
	TransferStatusRejected
	TransferStatusCancelled
)

type CertAlert struct {
	ID                    int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	SeafarerCertificateID int64     `json:"seafarer_certificate_id" gorm:"not null"`
	SeafarerID            int64     `json:"seafarer_id" gorm:"not null;index"`
	AlertLevel            int8      `json:"alert_level" gorm:"not null;index"`
	AlertDate             time.Time `json:"alert_date" gorm:"not null;index"`
	ExpireDate            time.Time `json:"expire_date" gorm:"not null"`
	DaysRemaining         int       `json:"days_remaining" gorm:"not null"`
	IsHandled             int8      `json:"is_handled" gorm:"not null;default:0;index"`
	HandleRemark          string    `json:"handle_remark" gorm:"size:500"`
	CreatedAt             time.Time `json:"created_at"`

	SeafarerCertificate *SeafarerCertificate `json:"certificate,omitempty" gorm:"foreignKey:SeafarerCertificateID"`
	Seafarer            *Seafarer            `json:"seafarer,omitempty" gorm:"foreignKey:SeafarerID"`
}

func (CertAlert) TableName() string {
	return "cert_alert"
}

const (
	AlertLevel90Days = iota + 1
	AlertLevel60Days
	AlertLevel30Days
)
