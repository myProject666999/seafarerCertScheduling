package model

import "time"

type CertificateType struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string    `json:"name" gorm:"size:100;not null"`
	Code            string    `json:"code" gorm:"size:50;not null;uniqueIndex"`
	Description     string    `json:"description" gorm:"type:text"`
	ValidityMonths  *int      `json:"validity_months"`
	IsRequired      int8      `json:"is_required" gorm:"not null;default:0"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (CertificateType) TableName() string {
	return "certificate_type"
}

type SeafarerCertificate struct {
	ID                int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	SeafarerID        int64          `json:"seafarer_id" gorm:"not null;index"`
	CertificateTypeID int64          `json:"certificate_type_id" gorm:"not null"`
	CertNumber        string         `json:"cert_number" gorm:"size:100;not null"`
	IssueDate         time.Time      `json:"issue_date" gorm:"not null"`
	ExpireDate        *time.Time     `json:"expire_date" gorm:"index"`
	CertImageURL      string         `json:"cert_image_url" gorm:"size:500"`
	Status            int8           `json:"status" gorm:"not null;default:1;index"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         *time.Time     `json:"deleted_at"`

	Seafarer        *Seafarer        `json:"seafarer,omitempty" gorm:"foreignKey:SeafarerID"`
	CertificateType *CertificateType `json:"certificate_type,omitempty" gorm:"foreignKey:CertificateTypeID"`
	Alerts          []CertAlert      `json:"alerts,omitempty" gorm:"foreignKey:SeafarerCertificateID"`
}

func (SeafarerCertificate) TableName() string {
	return "seafarer_certificate"
}

const (
	CertStatusExpired = iota
	CertStatusValid
	CertStatusExpiring
)
