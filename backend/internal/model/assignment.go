package model

import "time"

type SeafarerAssignment struct {
	ID                     int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	SeafarerID             int64      `json:"seafarer_id" gorm:"not null;index"`
	ShipID                 int64      `json:"ship_id" gorm:"not null;index"`
	ShipPositionID         int64      `json:"ship_position_id" gorm:"not null;index"`
	EmbarkDate             time.Time  `json:"embark_date" gorm:"not null"`
	ExpectedDisembarkDate  *time.Time `json:"expected_disembark_date"`
	ActualDisembarkDate    *time.Time `json:"actual_disembark_date"`
	Status                 int8       `json:"status" gorm:"not null;default:1;index"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`

	Seafarer     *Seafarer     `json:"seafarer,omitempty" gorm:"foreignKey:SeafarerID"`
	Ship         *Ship         `json:"ship,omitempty" gorm:"foreignKey:ShipID"`
	ShipPosition *ShipPosition `json:"ship_position,omitempty" gorm:"foreignKey:ShipPositionID"`
}

func (SeafarerAssignment) TableName() string {
	return "seafarer_assignment"
}

const (
	AssignmentStatusDisembarked = iota
	AssignmentStatusOnboard
)
