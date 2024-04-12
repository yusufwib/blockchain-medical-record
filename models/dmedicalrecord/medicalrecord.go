package dmedicalrecord

import "time"

type MedicalRecord struct {
	DoctorID               uint64    `json:"doctor_id"`
	PatientID              uint64    `json:"patient_id"`
	AppointmentID          uint64    `json:"appointment_id"`
	Diagnose               string    `json:"diagnose"`
	Notes                  string    `json:"notes"`
	AdditionalDocumentPath string    `json:"additional_document_path"`
	Prescription           string    `json:"prescription"`
	CreatedAt              time.Time `json:"created_at"`
	MedicalRecordNumber    string    `json:"medical_record_number"`
}
