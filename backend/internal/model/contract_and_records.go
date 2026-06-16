package model

import "time"

type VoyageContract struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	SeafarerID     int64      `json:"seafarer_id" gorm:"not null;index"`
	ShipID         int64      `json:"ship_id" gorm:"not null;index"`
	ContractNumber string     `json:"contract_number" gorm:"size:50;not null;index"`
	StartDate      time.Time  `json:"start_date" gorm:"not null"`
	EndDate        time.Time  `json:"end_date" gorm:"not null"`
	ActualEndDate  *time.Time `json:"actual_end_date"`
	Status         int8       `json:"status" gorm:"not null;default:1"`
	Remarks        string     `json:"remarks" gorm:"type:text"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Seafarer *Seafarer `json:"seafarer,omitempty" gorm:"foreignKey:SeafarerID"`
	Ship     *Ship     `json:"ship,omitempty" gorm:"foreignKey:ShipID"`
}

func (VoyageContract) TableName() string {
	return "voyage_contract"
}

const (
	ContractStatusTerminated = iota
	ContractStatusActive
	ContractStatusCompleted
)

type EmbarkDisembarkRecord struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	SeafarerID  int64     `json:"seafarer_id" gorm:"not null;index"`
	ShipID      int64     `json:"ship_id" gorm:"not null;index"`
	RecordType  int8      `json:"record_type" gorm:"not null"`
	RecordDate  time.Time `json:"record_date" gorm:"not null;index"`
	Port        string    `json:"port" gorm:"size:100"`
	Reason      string    `json:"reason" gorm:"size:200"`
	Operator    string    `json:"operator" gorm:"size:50"`
	CreatedAt   time.Time `json:"created_at"`

	Seafarer *Seafarer `json:"seafarer,omitempty" gorm:"foreignKey:SeafarerID"`
	Ship     *Ship     `json:"ship,omitempty" gorm:"foreignKey:ShipID"`
}

func (EmbarkDisembarkRecord) TableName() string {
	return "embark_disembark_record"
}

const (
	RecordTypeEmbark = iota + 1
	RecordTypeDisembark
)

type LeaveRecord struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	SeafarerID  int64      `json:"seafarer_id" gorm:"not null;index"`
	StartDate   time.Time  `json:"start_date" gorm:"not null"`
	EndDate     *time.Time `json:"end_date"`
	LeaveDays   *int       `json:"leave_days"`
	Status      int8       `json:"status" gorm:"not null;default:1;index"`
	Reason      string     `json:"reason" gorm:"size:200"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Seafarer *Seafarer `json:"seafarer,omitempty" gorm:"foreignKey:SeafarerID"`
}

func (LeaveRecord) TableName() string {
	return "leave_record"
}

const (
	LeaveStatusEnded = iota
	LeaveStatusActive
)

type HealthReexamination struct {
	ID              int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	SeafarerID      int64      `json:"seafarer_id" gorm:"not null;index"`
	ExamDate        time.Time  `json:"exam_date" gorm:"not null"`
	NextExamDate    *time.Time `json:"next_exam_date" gorm:"index"`
	ExamResult      int8       `json:"exam_result" gorm:"not null"`
	ExamInstitution string     `json:"exam_institution" gorm:"size:200"`
	ReportURL       string     `json:"report_url" gorm:"size:500"`
	Restrictions    string     `json:"restrictions" gorm:"type:text"`
	CreatedAt       time.Time  `json:"created_at"`

	Seafarer *Seafarer `json:"seafarer,omitempty" gorm:"foreignKey:SeafarerID"`
}

func (HealthReexamination) TableName() string {
	return "health_reexamination"
}

const (
	ExamResultQualified = iota + 1
	ExamResultUnqualified
	ExamResultRestricted
)
