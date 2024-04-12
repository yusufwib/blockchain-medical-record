package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yusufwib/blockchain-medical-record/config"
	"github.com/yusufwib/blockchain-medical-record/models/dappointment"
	"github.com/yusufwib/blockchain-medical-record/models/dmedicalrecordaccess"
	"gorm.io/gorm"
)

type AppointmentRepository struct {
	DB     *gorm.DB
	Config *config.ConfigGroup
}

func NewAppointmentRepository(DB *gorm.DB, cfg *config.ConfigGroup) AppointmentRepository {
	return AppointmentRepository{DB, cfg}
}

func (r *AppointmentRepository) session(ctx context.Context) *gorm.DB {
	trx, ok := ctx.Value("pg").(*gorm.DB)
	if !ok {
		return r.DB
	}
	return trx
}

func (r *AppointmentRepository) FindAppointmentByPatientID(ctx context.Context, ID uint64, filter dappointment.AppointmentFilter) (res []dappointment.AppointmentResponse, err error) {
	trx := r.session(ctx)
	ctxWT, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if filter.IsDoctor && filter.PatientID != 0 && filter.AppointmentID != 0 {
		var medicalrecordaccess dmedicalrecordaccess.MedicalRecordAccess
		if err = trx.Debug().WithContext(ctxWT).Table(dmedicalrecordaccess.TableName()).
			Joins("JOIN appointments ON medical_record_accesses.appointment_id = appointments.id").
			Where("appointments.id = ?", filter.AppointmentID).
			Where("status = ?", dappointment.AppointmentStatusUpcoming).
			First(&medicalrecordaccess).
			Error; err != nil {
			return nil, fmt.Errorf("error while retrieving access key: %w", err)
		}

		if err = r.checkAccessKey(medicalrecordaccess.AccessKey, fmt.Sprintf("%d", ID)); err != nil {
			return nil, fmt.Errorf("error while check access key: %w", err)
		}
	}

	query := trx.Debug().WithContext(ctxWT).Table(dappointment.TableName()).
		Select("appointments.*, u1.name AS doctor_name, u2.name AS patient_name, hs.name AS health_service_name, patients.allergies, appointments.created_at AS booking_at").
		Joins("JOIN doctors ON appointments.doctor_id = doctors.id").
		Joins("JOIN users u1 ON doctors.user_id = u1.id").
		Joins("JOIN patients ON appointments.patient_id = patients.id").
		Joins("JOIN users u2 ON patients.user_id = u2.id").
		Joins("JOIN health_services hs ON doctors.health_service_id = hs.id")

	if !filter.IsDoctor && filter.PatientID == 0 {
		query = query.Where("appointments.patient_id = ?", ID)
	} else {
		query = query.Where("appointments.doctor_id = ?", ID)
	}

	if filter.IsDoctor && filter.PatientID != 0 && filter.AppointmentID != 0 {
		query = query.Where("appointments.patient_id = ?", filter.PatientID)
	}

	if filter.Status != "" {
		query = query.Where("status IN (?)", strings.Split(filter.Status, ","))
	}

	if filter.HealthServiceID != 0 {
		query = query.Where("appointments.health_service_id = ?", filter.HealthServiceID)
	}

	if filter.ScheduleDate != "" {
		query = query.Where("schedule_date = ?", filter.ScheduleDate)
	}

	if err = query.Find(&res).Error; err != nil {
		return nil, fmt.Errorf("error while retrieving appointments: %w", err)
	}

	return
}

func (r *AppointmentRepository) FindAppointmentDetailByID(ctx context.Context, ID uint64) (res dappointment.AppointmentResponse, err error) {
	trx := r.session(ctx)
	ctxWT, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	query := trx.Debug().WithContext(ctxWT).Table(dappointment.TableName()).
		Select("appointments.*, u1.name AS doctor_name, u2.name AS patient_name, hs.name AS health_service_name, patients.allergies, appointments.created_at AS booking_at").
		Joins("JOIN doctors ON appointments.doctor_id = doctors.id").
		Joins("JOIN users u1 ON doctors.user_id = u1.id").
		Joins("JOIN patients ON appointments.patient_id = patients.id").
		Joins("JOIN users u2 ON patients.user_id = u2.id").
		Joins("JOIN health_services hs ON doctors.health_service_id = hs.id")

	if err = query.Where("appointments.id = ?", ID).First(&res).Error; err != nil {
		return dappointment.AppointmentResponse{}, fmt.Errorf("error while retrieving appointment detail: %w", err)
	}

	return
}

func (r *AppointmentRepository) CreateAppointment(ctx context.Context, patientID uint64, req dappointment.AppointmentCreateRequest) (err error) {
	trx := r.session(ctx)
	ctxWT, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	appointment := req.ToAppointment(patientID)
	if err = trx.WithContext(ctxWT).Table(dappointment.TableName()).Create(&appointment).Error; err != nil {
		return fmt.Errorf("error while create appointments: %w", err)
	}

	accessKey, err := r.generateAccessKey(appointment)
	if err != nil {
		return fmt.Errorf("error while generate access key: %w", err)
	}

	medicalRecordAccess := req.ToMedicalRecordAccess(appointment, accessKey)
	if err = trx.WithContext(ctxWT).Table(dmedicalrecordaccess.TableName()).Create(&medicalRecordAccess).Error; err != nil {
		return fmt.Errorf("error while create medical record access: %w", err)
	}

	return
}

func (r *AppointmentRepository) UpdateAppointmentStatus(ctx context.Context, ID uint64, req dappointment.AppointmentUpdateStatusRequest) (err error) {
	trx := r.session(ctx)
	ctxWT, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err = trx.WithContext(ctxWT).Table(dappointment.TableName()).
		Where("id = ?", ID).
		Updates(map[string]interface{}{
			"status": req.Status,
		}).Error; err != nil {
		return fmt.Errorf("error while update status appointments: %w", err)
	}

	return
}

func (r *AppointmentRepository) generateAccessKey(appointment dappointment.Appointment) (string, error) {
	claims := jwt.MapClaims{
		"appointment_id": appointment.ID,
		"patient_id":     appointment.PatientID,
		"doctor_id":      appointment.DoctorID,
		"schedule_date":  appointment.ScheduleDate,
		"schedule_time":  appointment.ScheduleTime,
		"status":         appointment.Status,
		"exp":            time.Now().Add(time.Hour * 24 * 30 * 365).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(fmt.Sprintf("%d", appointment.DoctorID)))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (r *AppointmentRepository) checkAccessKey(accessKey, privateKey string) (err error) {
	token, err := jwt.Parse(accessKey, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid private key")
		}
		return []byte(privateKey), nil
	})

	if err != nil {
		return fmt.Errorf("error while parsing access key: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid access key claims")
	}

	scheduleDate, _ := claims["schedule_date"].(string)
	scheduleTime, _ := claims["schedule_time"].(string)

	schedule, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", scheduleDate, scheduleTime))
	if err != nil {
		return fmt.Errorf("error while parsing schedule date or time: %w", err)
	}

	if schedule.Before(time.Now()) {
		return fmt.Errorf("invalid access key: schedule date or time is invalid")
	}

	return nil
}
