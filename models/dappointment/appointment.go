package dappointment

import (
	"time"
)

type AppointmentStatus string

const (
	AppointmentStatusWaiting   AppointmentStatus = "WAITING"
	AppointmentStatusUpcoming  AppointmentStatus = "UPCOMING"
	AppointmentStatusDone      AppointmentStatus = "DONE"
	AppointmentStatusCancelled AppointmentStatus = "CANCELLED"
	AppointmentStatusRejected  AppointmentStatus = "REJECTED"
)

type Appointment struct {
	ID              uint64            `gorm:"column:id;primaryKey" json:"id"`
	RecordNumber    string            `gorm:"column:record_number;not null" json:"record_number"`
	DoctorID        uint64            `gorm:"column:doctor_id;not null" json:"doctor_id"`
	PatientID       uint64            `gorm:"column:patient_id;not null" json:"patient_id"`
	HealthServiceID uint64            `gorm:"column:health_service_id;not null" json:"health_service_id"`
	Symptoms        string            `gorm:"column:symptoms" json:"symptoms"`
	Status          AppointmentStatus `gorm:"column:status;not null" json:"status"`
	ScheduleDate    string            `gorm:"column:schedule_date" json:"schedule_date"`
	ScheduleTime    string            `gorm:"column:schedule_time" json:"schedule_time"`
	CreatedAt       time.Time         `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func TableName() string {
	return "appointments"
}
