package dmedicalrecord

import "mime/multipart"

type UploadFileRequest struct {
	File *multipart.FileHeader
}

type MedicalRecordRequest struct {
	DoctorID               uint64 `json:"doctor_id" validate:"required"`
	PatientID              uint64 `json:"patient_id" validate:"required"`
	AppointmentID          uint64 `json:"appointment_id"`
	Diagnose               string `json:"diagnose" validate:"required"`
	Notes                  string `json:"notes" validate:"required"`
	AdditionalDocumentPath string `json:"additional_document_path" validate:"required"`
	Prescription           string `json:"prescription" validate:"required"`
}
