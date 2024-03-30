package service

import (
	"context"

	"github.com/yusufwib/blockchain-medical-record/models/dappointment"
	"github.com/yusufwib/blockchain-medical-record/repository"
)

type AppointmentService struct {
	AppointmentRepository repository.AppointmentRepository
}

func NewAppointmentService(r repository.AppointmentRepository) AppointmentService {
	return AppointmentService{
		AppointmentRepository: r,
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
