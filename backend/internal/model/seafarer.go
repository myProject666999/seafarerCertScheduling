package model

import "time"

type Seafarer struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string     `json:"name" gorm:"size:50;not null"`
	Gender    int8       `json:"gender" gorm:"not null;default:0"`
	Birthday  *time.Time `json:"birthday"`
	IDNumber  string     `json:"id_number" gorm:"size:18"`
	Phone     string     `json:"phone" gorm:"size:20"`
	Email     string     `json:"email" gorm:"size:100"`
	Rank      string     `json:"rank" gorm:"size:50"`
	Status    int8       `json:"status" gorm:"not null;default:0"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`

	Certificates          []SeafarerCertificate  `json:"certificates,omitempty" gorm:"foreignKey:SeafarerID"`
	Assignments           []SeafarerAssignment   `json:"assignments,omitempty" gorm:"foreignKey:SeafarerID"`
	VoyageContracts       []VoyageContract       `json:"voyage_contracts,omitempty" gorm:"foreignKey:SeafarerID"`
	EmbarkDisembarkRecords []EmbarkDisembarkRecord `json:"embark_records,omitempty" gorm:"foreignKey:SeafarerID"`
	LeaveRecords          []LeaveRecord          `json:"leave_records,omitempty" gorm:"foreignKey:SeafarerID"`
	HealthReexaminations  []HealthReexamination  `json:"health_reexams,omitempty" gorm:"foreignKey:SeafarerID"`
}

func (Seafarer) TableName() string {
	return "seafarer"
}

const (
	SeafarerStatusPending = iota
	SeafarerStatusOnboard
	SeafarerStatusLeave
)
