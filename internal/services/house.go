package services

import (
	"avito/internal/domain"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type HouseRepository interface {
	Create(ctx context.Context, house domain.House) (domain.House, error)
	GetByID(ctx context.Context, houseID int) (domain.House, error)
	GetAll(ctx context.Context, offset int, limit int) ([]domain.House, error)
	GetFlatsByHouseID(ctx context.Context, houseID int) ([]domain.Flat, error)
	SubscribeByID(ctx context.Context, houseID int, userID uuid.UUID) error
}

type HouseServce struct {
	houseRepo HouseRepository
	log       *slog.Logger
	// notifySender domain.NotifySender
	// notifyRepo   domain.NotifyRepo
}

func NewHouseServce(
	houseRepo HouseRepository,
	// notifySender domain.NotifySender,
	// notifyRepo domain.NotifyRepo,
	logger *slog.Logger,
	done chan bool,
	freq time.Duration,
	timeout time.Duration,
) *HouseServce {
	houseUsecase := HouseServce{
		houseRepo: houseRepo,
		log:       logger,
		// notifySender: notifySender,
		// notifyRepo:   notifyRepo,
	}

	// go houseUsecase.Notifying(done, freq, timeout)

	return &houseUsecase
}

func (s *HouseServce) Create(ctx context.Context, houseRequest domain.HouseCreateRequest) (domain.HouseCreateResponse, error) {
	const op = "HouseServce.Create"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("creating house")

	date := time.Now()
	house := domain.House{
		Address:   houseRequest.Address,
		Year:      houseRequest.Year,
		Developer: houseRequest.Developer,
		CreatedAt: date,
		UpdatedAt: date,
	}

	house, err := s.houseRepo.Create(ctx, house)
	if err != nil {
		s.log.Warn("create house error: " + err.Error())
		return domain.HouseCreateResponse{},
			fmt.Errorf("user usecase: create error: %v", err.Error())
	}

	houseResponse := domain.HouseCreateResponse{
		HomeID:    house.HouseID,
		Address:   house.Address,
		Year:      house.Year,
		Developer: house.Developer,
		CreatedAt: house.CreatedAt,
		UpdateAt:  house.UpdatedAt,
	}

	return houseResponse, nil
}

func (s *HouseServce) GetFlatsByHouseID(ctx context.Context, houseID int) (domain.HouseFlatsResponse, error) {
	const op = "HouseServce.GetFlatsByHouseID"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("getting flats by house_id")

	flats, err := s.houseRepo.GetFlatsByHouseID(ctx, houseID)
	if err != nil {
		s.log.Warn("get flats by house_id error: " + err.Error())
		return domain.HouseFlatsResponse{}, fmt.Errorf("get flats by house_id error: %v", err.Error())
	}

	var flatsResponse domain.HouseFlatsResponse
	for i := 0; i < len(flats); i++ {
		var tmpFlat domain.HouseFlatResponse
		tmpFlat.HouseID = flats[i].HouseID
		tmpFlat.Price = flats[i].Price
		tmpFlat.Rooms = flats[i].Rooms
		tmpFlat.Status = flats[i].Status
		flatsResponse.Flats = append(flatsResponse.Flats, tmpFlat)
	}

	return flatsResponse, nil
}

func (s *HouseServce) SubscribeByID(ctx context.Context, houseID int, userID uuid.UUID) error {
	const op = "HouseServce.SubscribeByID"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("subscribing by id")

	err := s.houseRepo.SubscribeByID(ctx, houseID, userID)
	if err != nil {
		s.log.Warn("subscribing by id: " + err.Error())
		return fmt.Errorf("subscribing by id: %v", err.Error())
	}

	return nil
}

// func (s *HouseServce) Notifying(done chan bool, frequency time.Duration, timeout time.Duration, lg *zap.Logger) {
// 	for {
// 		select {
// 		case <-done:
// 			lg.Warn("house usecase: subscribing goroutine exited")
// 			return
// 		default:
// 			lg.Info("house usecase: subscribing goroutine working")
// 			ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 			defer cancel()

// 			notifies, err := s.notifyRepo.GetNoSendNotifies(ctx, lg)
// 			if err != nil {
// 				lg.Warn("house usecase: notifying error", zap.Error(err))
// 			}

// 			for _, notify := range notifies {
// 				msg := fmt.Sprintf("New flat with number %d in house %d!", notify.FlatID, notify.HouseID)
// 				err = s.notifySender.SendEmail(ctx, notify.UserMail, msg)
// 				if err != nil {
// 					lg.Warn("house usecase: notifying error: send email error", zap.Error(err))
// 					continue
// 				} else {
// 					err = s.notifyRepo.SendNotifyByID(ctx, notify.ID, lg)
// 				}
// 			}
// 			time.Sleep(frequency)
// 		}
// 	}
// }
