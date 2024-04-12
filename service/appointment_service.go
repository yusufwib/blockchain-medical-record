package service

import (
	"context"
	"io"
	"os"
	"path/filepath"

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

func (s AppointmentService) CreateAppointment(ctx context.Context, patientID uint64, req dappointment.AppointmentCreateRequest) error {
	return s.AppointmentRepository.CreateAppointment(ctx, patientID, req)
}

func (s AppointmentService) UpdateAppointmentStatus(ctx context.Context, ID uint64, req dappointment.AppointmentUpdateStatusRequest) (err error) {
	return s.AppointmentRepository.UpdateAppointmentStatus(ctx, ID, req)
}

func (s AppointmentService) UploadFile(ctx context.Context, req dmedicalrecord.UploadFileRequest) (path string, err error) {
	src, err := req.File.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dstPath := filepath.Join("public", "documents", req.File.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
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
