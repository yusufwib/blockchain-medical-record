package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
	"github.com/yusufwib/blockchain-medical-record/models/dappointment"
	"github.com/yusufwib/blockchain-medical-record/models/dmedicalrecord"
	"github.com/yusufwib/blockchain-medical-record/repository"
	mlog "github.com/yusufwib/blockchain-medical-record/utils/logger"
)

type AppointmentService struct {
	AppointmentRepository repository.AppointmentRepository
	BlockchainRepository  repository.BlockchainRepository
	Logger                mlog.Logger
}

func NewAppointmentService(ra repository.AppointmentRepository, rb repository.BlockchainRepository, logger mlog.Logger) AppointmentService {
	return AppointmentService{
		AppointmentRepository: ra,
		BlockchainRepository:  rb,
		Logger:                logger,
	}
}

func (s AppointmentService) FindAppointmentByPatientID(ctx context.Context, ID uint64, filter dappointment.AppointmentFilter) ([]dappointment.AppointmentResponse, error) {
	return s.AppointmentRepository.FindAppointmentByPatientID(ctx, ID, filter)
}

func (s AppointmentService) FindAppointmentDetailByID(ctx context.Context, ID uint64) (dappointment.AppointmentResponse, error) {
	return s.AppointmentRepository.FindAppointmentDetailByID(ctx, ID)
}

func (s AppointmentService) CreateAppointment(ctx context.Context, patientID uint64, req dappointment.AppointmentCreateRequest) (uint64, error) {
	return s.AppointmentRepository.CreateAppointment(ctx, patientID, req)
}

func (s AppointmentService) UpdateAppointmentStatus(ctx context.Context, ID uint64, req dappointment.AppointmentUpdateStatusRequest) (err error) {
	return s.AppointmentRepository.UpdateAppointmentStatus(ctx, ID, req)
}

func (s AppointmentService) UploadFile(ctx context.Context, req dmedicalrecord.UploadFileRequest) (path string, err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	src, err := req.File.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dstPath := filepath.Join("public", "documents", req.File.Filename)
	realDstPath := filepath.Join(currentDir, "public/documents", req.File.Filename)

	dst, err := os.Create(realDstPath)
	if err != nil {
		fmt.Println("error create file", err)
		return "error", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		fmt.Println("error copy file", err)
		return "", err
	}

	return dstPath, nil
}

func (s AppointmentService) WriteMedicalRecord(ctx context.Context, req dmedicalrecord.MedicalRecordRequest) (err error) {
	s.BlockchainRepository.AddBlockMedicalRecord(req)

	return s.AppointmentRepository.UpdateAppointmentStatus(ctx, req.AppointmentID, dappointment.AppointmentUpdateStatusRequest{
		Status: dappointment.AppointmentStatusDone,
	})
}

func (s AppointmentService) FindMedicalRecordByID(ctx context.Context, ID uint64) (res dmedicalrecord.MedicalRecord) {
	return s.BlockchainRepository.GetBlocksByAppointmentID(ID)
}

func (s AppointmentService) ExportMedicalRecord(ctx context.Context, ID uint64) (res string, err error) {
	// appointment, err := s.AppointmentRepository.FindAppointmentDetailByID(ctx, ID)
	// if err != nil {
	// 	return
	// }

	// if appointment.IsEmpty() {
	// 	return "", fmt.Errorf("appointment is empty")
	// }

	// medcialRecord := s.BlockchainRepository.GetBlocksByAppointmentID(ID)

	// if medcialRecord.IsEmpty() {
	// 	return "", fmt.Errorf("medical record is empty")
	// }
	// Create new PDF document
	// Create new PDF document
	// Create new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Add a new page to the PDF
	pdf.AddPage()

	// Set font for keys (labels)
	pdf.SetFont("Arial", "B", 16)

	// Add title
	pdf.Cell(0, 10, "Rekam Medis Elektronik")
	pdf.Ln(12) // Line break

	// Define content
	content1 := map[string]string{
		"Pasien:": "Hafiz Vario",
		"Dokter:": "Dr. Erina Spesialis Tendangan Bebas",
	}

	// Add content to PDF
	for label, value := range content1 {
		pdf.Ln(8)
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(40, 10, label)
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(40, 10, value)
	}
	pdf.Ln(17)
	pdf.SetFont("Arial", "B", 12)

	// Add title
	pdf.Cell(0, 15, "Detail Rekam Medis")
	pdf.Ln(4)
	// Define content
	content := map[string]string{
		"ID Rekam Medis:":    "MR-202408071206",
		"Tanggal Penulisan:": "28 Maret 2024",
		"Pemeriksa:":         "Dr. Yusuf Banana - Dokter Umum",
		"Usia:":              "22",
		"Berat Badan:":       "75 kg",
		"Tinggi Badan:":      "175 cm",
		"Alergi:":            "Obat Paracetamol & Udang",
		"Keluhan:":           "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"Catatan:":           "Lorem ipsum",
		"Gambar:":            "...",
	}

	// Add content to PDF
	for label, value := range content {
		pdf.Ln(8)
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(40, 10, label)
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(40, 10, value)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// Save PDF to file

	path := fmt.Sprintf("public/medical-record/%d.pdf", ID)
	err = pdf.OutputFileAndClose(currentDir + "/" + path)
	if err != nil {
		// panic(err)
		return res, err
	}

	return path, nil
}
