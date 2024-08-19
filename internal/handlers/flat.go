package handlers

import (
	"avito/internal/domain"
	manager "avito/pkg/jwt"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type FlatService interface {
	Create(ctx context.Context, flatRequest domain.FlatCreateRequest) (domain.FlatCreateResponse, error)
	Update(ctx context.Context, moderatorID uuid.UUID, newFlatData domain.FlatUpdateRequest) (domain.FlatCreateResponse, error)
}

type FlatHandler struct {
	flatService FlatService
	manager     manager.TokenManager
	log         *slog.Logger
}

func NewFlatHandler(flatService FlatService, manager manager.TokenManager, log *slog.Logger) *FlatHandler {
	return &FlatHandler{
		flatService: flatService,
		manager:     manager,
		log:         log,
	}
}
