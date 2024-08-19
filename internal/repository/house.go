package repository

import (
	"avito/internal/domain"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HouseRepository struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewPostgresHouseRepo(pool *pgxpool.Pool, logger *slog.Logger) *HouseRepository {
	return &HouseRepository{
		pool: pool,
		log:  logger,
	}
}

func (r *HouseRepository) Create(ctx context.Context, house domain.House) (domain.House, error) {
	r.log.Info("create house", slog.Int("house_id", house.HouseID))

	var savedHouse domain.House

	query := `insert into houses(address, year, developer, created_at, updated_at)
	values ($1, $2, $3, $4, $5) returning house_id`
	rows := r.pool.QueryRow(ctx, query,
		house.Address,
		house.Year,
		house.Developer,
		house.CreatedAt,
		house.UpdatedAt)
	err := rows.Scan(
		&savedHouse.HouseID,
		&savedHouse.Address,
		&savedHouse.Year,
		&savedHouse.Developer,
		&savedHouse.CreatedAt,
		&savedHouse.UpdatedAt)
	if err != nil {
		r.log.Warn("house create error: " + err.Error())
		return domain.House{}, err
	}

	return savedHouse, nil
}

func (r *HouseRepository) GetByID(ctx context.Context, houseID int) (domain.House, error) {
	r.log.Info("get house by id", slog.Int("id", houseID))

	var house domain.House

	query := `select house_id, address, year, developer, created_at, updated_at from houses where house_id = $1`
	rows := r.pool.QueryRow(ctx, query, houseID)
	err := rows.Scan(
		&house.HouseID,
		&house.Address,
		&house.Year,
		&house.Developer,
		&house.CreatedAt,
		&house.UpdatedAt)
	if err != nil {
		r.log.Warn("house get by id error: " + err.Error())
		return domain.House{}, err
	}

	return house, nil
}

func (r *HouseRepository) GetAll(ctx context.Context, offset int, limit int) ([]domain.House, error) {
	r.log.Info("get all houses", slog.Int("offset", offset), slog.Int("limit", limit))

	query := `select house_id, address, year, developer, created_at, updated_at from houses limit $1 offset $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		r.log.Warn("get all error: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var (
		houses []domain.House
		house  domain.House
	)
	for rows.Next() {
		err = rows.Scan(
			&house.HouseID,
			&house.Address,
			&house.Year,
			&house.Developer,
			&house.CreatedAt,
			&house.UpdatedAt)
		if err != nil {
			r.log.Warn("scan house error")
			continue
		}
		houses = append(houses, house)
	}

	return houses, err
}

func (r *HouseRepository) GetFlatsByHouseID(ctx context.Context, houseID int) ([]domain.Flat, error) {
	r.log.Info("get flats by house_id", slog.Int("house_id", houseID))

	query := `select flats.flat_id, houses.house_id, flats.price, flats.rooms, flats.status 
			from flats join houses
			on flats.house_id = houses.house_id
			where houses.house_id = $1`

	rows, err := r.pool.Query(ctx, query, houseID)
	if err != nil {
		r.log.Warn("get flats by house_id failed: " + err.Error())
		return nil, fmt.Errorf("get flats by house_id: %v", err.Error())
	}
	defer rows.Close()

	var flats []domain.Flat
	for rows.Next() {
		flat := domain.Flat{}
		err = rows.Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Rooms, &flat.Status)
		if err != nil {
			r.log.Warn("scan house error: " + err.Error())
			continue
		}
		flats = append(flats, flat)
	}

	return flats, err
}

func (r *HouseRepository) SubscribeByID(ctx context.Context, houseID int, userID uuid.UUID) error {
	r.log.Info("dubscribe by house_id")

	query := `insert into subscribers(user_id, house_id) values ($1, $2)`
	_, err := r.pool.Exec(ctx, query, userID, houseID)
	if err != nil {
		r.log.Warn("subscribe by house_id error: " + err.Error())
		return fmt.Errorf("subscribe by id error: %v", err.Error())
	}

	return nil
}
