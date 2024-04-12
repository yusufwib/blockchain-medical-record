package dappointment

import (
	"time"

	"github.com/yusufwib/blockchain-medical-record/models/dmedicalrecordaccess"
	"github.com/yusufwib/blockchain-medical-record/utils/randstr"
)

type AppointmentCreateRequest struct {
	ID              uint64 `gorm:"column:id;primaryKey"`
	DoctorID        uint64 `json:"doctor_id" validate:"required"`
	HealthServiceID uint64 `json:"health_service_id" validate:"required"`
	Symptoms        string `json:"symptoms" validate:"required"`
	ScheduleDate    string `json:"schedule_date" validate:"required"`
	ScheduleTime    string `json:"schedule_time" validate:"required"`
}

type AppointmentUpdateStatusRequest struct {
	Status AppointmentStatus `json:"status" validate:"required"`
}

type AppointmentFilter struct {
	// general
	ScheduleDate    string `json:"schedule_date"`
	HealthServiceID uint64 `json:"health_service_id"`
	Status          string `json:"status"`

	// doctor
	IsDoctor      bool   `json:"is_doctor"`
	PatientID     uint64 `json:"patient_id"`
	AppointmentID uint64 `json:"appointment_id"`
}

func (a AppointmentCreateRequest) ToAppointment(patientID uint64) Appointment {
	return Appointment{
		RecordNumber:    randstr.GenerateRandomString("APPT"),
		DoctorID:        a.DoctorID,
		PatientID:       patientID,
		Status:          AppointmentStatusWaiting,
		Symptoms:        a.Symptoms,
		ScheduleDate:    a.ScheduleDate,
		ScheduleTime:    a.ScheduleTime,
		HealthServiceID: a.HealthServiceID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func (a AppointmentCreateRequest) ToMedicalRecordAccess(appointment Appointment, accessKey string) dmedicalrecordaccess.MedicalRecordAccess {
	return dmedicalrecordaccess.MedicalRecordAccess{
		AppointmentID: appointment.ID,
		DoctorID:      appointment.DoctorID,
		PatientID:     appointment.PatientID,
		AccessKey:     accessKey,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
