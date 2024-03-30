package dmedicalrecordaccess

import (
	"time"
)

type MedicalRecordAccess struct {
	ID            uint64    `gorm:"column:id;primaryKey" json:"id"`
	AppointmentID uint64    `gorm:"column:appointment_id;not null" json:"appointment_id"`
	DoctorID      uint64    `gorm:"column:doctor_id;not null" json:"doctor_id"`
	PatientID     uint64    `gorm:"column:patient_id;not null" json:"patient_id"`
	AccessKey     string    `gorm:"column:access_key" json:"access_key"`
	CreatedAt     time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func TableName() string {
	return "medical_record_accesses"
}
