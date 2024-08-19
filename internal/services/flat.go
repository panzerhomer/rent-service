package services

import (
	"avito/internal/domain"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type FlatRepository interface {
	Create(ctx context.Context, flat domain.Flat) (domain.Flat, error)
	Update(ctx context.Context, newFlat domain.Flat) (domain.Flat, error)
	GetByID(ctx context.Context, flatID int, houseID int) (domain.Flat, error)
	GetAll(ctx context.Context, offset int, limit int) ([]domain.Flat, error)
}

type FlatService struct {
	flatRepo FlatRepository
	log      *slog.Logger
}

func NewFlatService(flatRepo FlatRepository, log *slog.Logger) *FlatService {
	return &FlatService{
		flatRepo: flatRepo,
		log:      log,
	}
}

func IsCorrectStatus(status string) bool {
	return status == domain.StatusCreated || status == domain.StatusApproved ||
		status == domain.StatusOnModeration || status == domain.StatusDeclined
}

func (s *FlatService) Create(ctx context.Context, flatRequest domain.FlatCreateRequest) (domain.FlatCreateResponse, error) {
	const op = "FlatService.Create"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("creating flat")

	flat := domain.Flat{
		HouseID: flatRequest.HouseID,
		Price:   flatRequest.Price,
		Rooms:   flatRequest.Rooms,
		Status:  domain.StatusCreated,
	}

	createdFlat, err := s.flatRepo.Create(ctx, flat)
	if err != nil {
		s.log.Warn("creating flat error: " + err.Error())
		return domain.FlatCreateResponse{},
			fmt.Errorf("creating flat error: %v", err.Error())
	}

	createdFlatResponse := domain.FlatCreateResponse{
		ID:      createdFlat.ID,
		HouseID: createdFlat.HouseID,
		Price:   createdFlat.Price,
		Rooms:   createdFlat.Rooms,
		Status:  createdFlat.Status,
	}

	return createdFlatResponse, nil
}

func (s *FlatService) Update(ctx context.Context, moderatorID uuid.UUID, newFlatData domain.FlatUpdateRequest) (domain.FlatCreateResponse, error) {
	const op = "FlatService.Update"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("updating flat")

	flat := domain.Flat{
		ID:     newFlatData.ID,
		Status: newFlatData.Status,
	}

	updatedFlat, err := s.flatRepo.Update(ctx, flat)
	if err != nil {
		s.log.Warn("updating error: " + err.Error())
		return domain.FlatCreateResponse{},
			fmt.Errorf("updating error: %v", err.Error())
	}

	updatedFlatResponse := domain.FlatCreateResponse{
		ID:      updatedFlat.ID,
		HouseID: updatedFlat.HouseID,
		Price:   updatedFlat.Price,
		Rooms:   updatedFlat.Rooms,
		Status:  updatedFlat.Status,
	}

	return updatedFlatResponse, nil
}
