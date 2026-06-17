package model

import "time"

type Ship struct {
	ID           int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string     `json:"name" gorm:"size:100;not null"`
	IMONumber    string     `json:"imo_number" gorm:"size:20;uniqueIndex"`
	MMSI         string     `json:"mmsi" gorm:"size:20"`
	ShipType     string     `json:"ship_type" gorm:"size:50"`
	GrossTonnage *float64   `json:"gross_tonnage" gorm:"type:decimal(10,2)"`
	FlagState    string     `json:"flag_state" gorm:"size:50"`
	Status       int8       `json:"status" gorm:"not null;default:1"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`

	Positions   []ShipPosition        `json:"positions,omitempty" gorm:"foreignKey:ShipID"`
	Assignments []SeafarerAssignment  `json:"assignments,omitempty" gorm:"foreignKey:ShipID"`
}

func (Ship) TableName() string {
	return "ship"
}

const (
	ShipStatusScrapped = iota
	ShipStatusOperating
	ShipStatusMaintenance
)

type ShipPosition struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ShipID        int64     `json:"ship_id" gorm:"not null;index"`
	PositionName  string    `json:"position_name" gorm:"size:50;not null"`
	Department    string    `json:"department" gorm:"size:50"`
	RequiredCount int       `json:"required_count" gorm:"not null;default:1"`
	SortOrder     int       `json:"sort_order" gorm:"not null;default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Ship               *Ship                    `json:"ship,omitempty" gorm:"foreignKey:ShipID"`
	CertRequirements   []PositionCertRequirement `json:"cert_requirements,omitempty" gorm:"foreignKey:ShipPositionID"`
	Assignments        []SeafarerAssignment     `json:"assignments,omitempty" gorm:"foreignKey:ShipPositionID"`
}

func (ShipPosition) TableName() string {
	return "ship_position"
}

type PositionCertRequirement struct {
	ID                int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ShipPositionID    int64     `json:"ship_position_id" gorm:"not null"`
	CertificateTypeID int64     `json:"certificate_type_id" gorm:"not null"`
	IsMandatory       int8      `json:"is_mandatory" gorm:"not null;default:1"`
	CreatedAt         time.Time `json:"created_at"`

	ShipPosition    *ShipPosition    `json:"ship_position,omitempty" gorm:"foreignKey:ShipPositionID"`
	CertificateType *CertificateType `json:"certificate_type,omitempty" gorm:"foreignKey:CertificateTypeID"`
}

func (PositionCertRequirement) TableName() string {
	return "position_cert_requirement"
}
