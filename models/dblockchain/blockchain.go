package dblockchain

import (
	"github.com/yusufwib/blockchain-medical-record/models/dmedicalrecord"
)

// Block represents a single block in the blockchain.
type Block struct {
	Index         uint64                       `json:"index"`
	Timestamp     string                       `json:"timestamp"`
	PatientID     uint64                       `json:"patient_id"`
	DoctorID      uint64                       `json:"doctor_id"`
	AppointmentID uint64                       `json:"appointment_id"`
	EncryptedData string                       `json:"encrypted_data"` //dmedicalrecord.MedicalRecord
	Message       string                       `json:"message"`
	PrevHash      string                       `json:"prev_hash"`
	Hash          string                       `json:"hash"`
	Data          dmedicalrecord.MedicalRecord `json:"-"`
}
