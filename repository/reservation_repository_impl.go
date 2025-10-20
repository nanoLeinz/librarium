package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/helper"
	"github.com/nanoLeinz/librarium/internal/enum"
	"github.com/nanoLeinz/librarium/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ReservationRepositoryImpl struct {
	db  *gorm.DB
	log *log.Logger
}

func NewReservationRepository(db *gorm.DB, log *log.Logger) ReservationRepository {
	return &ReservationRepositoryImpl{
		db:  db,
		log: log,
	}
}

func (s *ReservationRepositoryImpl) logWithCtx(ctx context.Context, function string) *log.Entry {

	traceID := ctx.Value(helper.KeyCon("traceID"))
	traceID = traceID.(string)

	logger := s.log.WithFields(log.Fields{
		"traceID":  traceID,
		"function": function,
	})

	return logger

}
func (s *ReservationRepositoryImpl) Create(ctx context.Context, reservation *model.Reservation) (*model.Reservation, error) {
	logger := s.logWithCtx(ctx, "ReservationRepository.Create").
		WithFields(log.Fields{
			"bookID":   reservation.BookID,
			"memberID": reservation.MemberID,
			"status":   reservation.Status,
		})

	logger.Info("executing reservation insert query")

	result := s.db.WithContext(ctx).Exec("INSERT INTO reservations (book_id,member_id,status,queue_position,reservation_date,created_at) VALUES (?,?,?,?,?,?)",
		reservation.BookID.String(),
		reservation.MemberID.String(),
		reservation.Status,
		reservation.QueuePosition,
		reservation.ReservationDate,
		reservation.ReservationDate)

	err := result.Error
	if err != nil {
		logger.WithError(err).Error("failed to insert reservation")
		return nil, err
	} else if result.RowsAffected < 1 {
		logger.Warn("reservation insert query executed but no rows affected")
	} else {
		logger.Info("reservation inserted successfully")
	}

	return reservation, nil
}

func (s *ReservationRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Reservation, error) {
	logger := s.logWithCtx(ctx, "ReservationRepository.GetByID").
		WithField("reservationID", id)

	logger.Info("executing get reservation by ID query")

	resv := model.Reservation{}
	result := s.db.WithContext(ctx).Raw("SELECT * FROM reservations WHERE id = ?", id).Scan(&resv)
	if result.Error != nil {
		logger.WithError(result.Error).Error("failed to get reservation by ID")
		return nil, result.Error
	} else if result.RowsAffected < 1 {
		logger.Warn("no reservation found for given ID")
		return nil, gorm.ErrRecordNotFound
	}

	logger.Info("reservation fetched successfully")
	return &resv, nil
}

func (s *ReservationRepositoryImpl) Update(ctx context.Context, reservation model.Reservation) error {
	logger := s.logWithCtx(ctx, "ReservationRepository.Update").
		WithField("reservationID", reservation.ID)

	logger.Info("executing reservation update query")

	result := s.db.WithContext(ctx).
		Exec("UPDATE reservations SET status = ?, queue_position = ?, updated_at = ? WHERE id = ? ",
			reservation.Status,
			reservation.QueuePosition,
			time.Now().Local(),
			reservation.ID)

	if result.Error != nil {
		logger.WithError(result.Error).Error("failed to update reservation")
		return result.Error
	} else if result.RowsAffected < 1 {
		logger.Warn("reservation update query executed but no rows affected")
	} else {
		logger.Info("reservation updated successfully")
	}

	return nil
}

func (s *ReservationRepositoryImpl) DeleteById(ctx context.Context, id uuid.UUID) error {
	logger := s.logWithCtx(ctx, "ReservationRepository.DeleteById").
		WithField("reservationID", id)

	logger.Info("executing reservation delete query")

	result := s.db.WithContext(ctx).
		Exec("UPDATE reservations SET deleted_at = ? WHERE id = ? ",
			time.Now().Local(),
			id)

	if result.Error != nil {
		logger.WithError(result.Error).Error("failed to delete reservation")
		return result.Error
	} else if result.RowsAffected < 1 {
		logger.Warn("reservation delete query executed but no rows affected")
	} else {
		logger.Info("reservation deleted successfully")
	}

	return nil

}
func (s *ReservationRepositoryImpl) GetAll(ctx context.Context) ([]model.Reservation, error) {
	logger := s.logWithCtx(ctx, "ReservationRepository.GetAll")
	logger.Info("executing get all reservations query")

	resv := []model.Reservation{}
	result := s.db.WithContext(ctx).Raw("SELECT * FROM reservations").Scan(&resv)
	if result.Error != nil {
		logger.WithError(result.Error).Error("failed to get all reservations")
		return nil, result.Error
	} else if result.RowsAffected < 1 {
		logger.Warn("no reservations found")
		return nil, gorm.ErrRecordNotFound
	}

	logger.WithField("count", len(resv)).Info("all reservations fetched successfully")
	return resv, nil
}

func (s *ReservationRepositoryImpl) GetLastQueue(ctx context.Context, bookID uuid.UUID) int {
	logger := s.logWithCtx(ctx, "ReservationRepository.GetLatestQueue").
		WithField("bookID", bookID.String())

	logger.Info("executing get latest queue query")

	var last int
	result := s.db.WithContext(ctx).
		Raw("SELECT COALESCE(MAX(queue_position), 0) FROM reservations WHERE book_id = ?", bookID.String()).
		Scan(&last)

	if result.Error != nil {
		logger.WithError(result.Error).Error("failed to fetch latest queue position")
		return 0
	}

	logger.WithField("lastQueue", last).Info("latest queue position fetched successfully")
	return last
}

func (s *ReservationRepositoryImpl) UpdateRelatedQueue(ctx context.Context, bookID uuid.UUID) error {
	logger := s.logWithCtx(ctx, "ReservationRepository.UpdateRelatedQueue").
		WithField("bookID", bookID.String())

	logger.Info("executing update related queue query")

	result := s.db.WithContext(ctx).
		Exec("UPDATE reservations SET queue_position = queue_position - 1, update_at = ? WHERE book_id = ? and status != ?",
			time.Now().Local(),
			bookID.String(),
			enum.PendingReserv.String(),
		)

	if result.Error != nil {
		logger.WithError(result.Error).Error("failed to update related queue")
		return result.Error
	}

	logger.WithField("rowsAffected", result.RowsAffected).Info("update related queue executed")
	if result.RowsAffected < 1 {
		logger.WithField("rowsAffected", result.RowsAffected).Warn("update affected no rows")
	} else {
		logger.WithField("rowsAffected", result.RowsAffected).Info("related queue updated successfully")
	}

	return nil
}
