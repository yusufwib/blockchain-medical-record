package dblockchain

import (
	"github.com/yusufwib/blockchain-medical-record/models/dmedicalrecord"
)

// Block represents a single block in the blockchain.
type Block struct {
	Index         uint64
	Timestamp     string
	PatientID     uint64
	DoctorID      uint64
	AppointmentID uint64
	EncryptedData string //dmedicalrecord.MedicalRecord
	Message       string
	PrevHash      string
	Hash          string
	Data          dmedicalrecord.MedicalRecord
}
