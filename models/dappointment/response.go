package dappointment

type AppointmentResponse struct {
	ID                uint64            `gorm:"column:id;primaryKey" json:"id"`
	RecordNumber      string            `json:"record_number"`
	DoctorID          uint64            `json:"doctor_id"`
	DoctorName        string            `json:"doctor_name"`
	PatientID         uint64            `json:"patient_id"`
	PatientName       string            `json:"patient_name"`
	Symptoms          string            `json:"symptoms"`
	Status            AppointmentStatus `json:"status"`
	HealthServiceID   uint64            `json:"health_service_id"`
	HealthServiceName string            `json:"health_service_name"`
	ScheduleDate      string            `json:"schedule_date"`
	ScheduleTime      string            `json:"schedule_time"`
}
