package repository

import (
	"avito/internal/domain"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FlatRepository struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewFlatRepository(pool *pgxpool.Pool, log *slog.Logger) *FlatRepository {
	return &FlatRepository{
		pool: pool,
		log:  log,
	}
}

func (r *FlatRepository) Create(ctx context.Context, flat domain.Flat) (domain.Flat, error) {
	r.log.Info("creating flat")

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.log.Warn("creating flat error: " + err.Error())
		return domain.Flat{}, fmt.Errorf("creating flat error: %v", err.Error())
	}
	defer tx.Rollback(ctx)

	var savedFlat domain.Flat
	query := `insert into flats(house_id, price, rooms, status)
			values ($1, $2, $3, $4) 
			returning flat_id`
	err = tx.QueryRow(ctx, query,
		flat.HouseID,
		flat.Price,
		flat.Rooms,
		domain.StatusCreated).
		Scan(&savedFlat.ID)
	if err != nil {
		r.log.Warn("creating error: " + err.Error())
		return domain.Flat{}, fmt.Errorf("creating error: " + err.Error())
	}

	savedFlat.Status = domain.StatusCreated
	savedFlat.Rooms = flat.Rooms
	savedFlat.Price = flat.Price
	savedFlat.HouseID = flat.HouseID

	return savedFlat, nil
}

func (r *FlatRepository) Update(ctx context.Context, newFlat domain.Flat) (domain.Flat, error) {
	r.log.Info("updating flat")

	var flat domain.Flat
	query := `update flats set status = $1 where flat_id = $2
		returning flat_id, status`
	rows := r.pool.QueryRow(ctx, query, newFlat.Status, newFlat.ID)
	err := rows.Scan(&flat.ID, &flat.HouseID, &flat.UserID,
		&flat.Price, &flat.Rooms, &flat.Status)
	if err != nil {
		r.log.Warn("updating error: " + err.Error())
		return domain.Flat{}, fmt.Errorf("updating error: %v", err.Error())
	}

	return flat, nil
}

func (r *FlatRepository) GetByID(ctx context.Context, flatID int, houseID int) (domain.Flat, error) {
	r.log.Info("get by flat_id, house_id")

	query := `select flat_id, house_id, price, rooms, status
				from flats where flat_id = $1 and house_id = $2`
	var flat domain.Flat
	rows := r.pool.QueryRow(ctx, query, flatID, houseID)
	err := rows.Scan(
		&flat.ID,
		&flat.HouseID,
		&flat.Price,
		&flat.Rooms,
		&flat.Status)
	if err != nil {
		r.log.Warn("get flat by ids: " + err.Error())
		return domain.Flat{}, fmt.Errorf("get flat by ids: %v", err.Error())
	}

	return flat, nil
}

func (r *FlatRepository) GetAll(ctx context.Context, offset int, limit int) ([]domain.Flat, error) {
	r.log.Info("getting all flats")

	query := `select flat_id, house_id, price, rooms, status from flats limit $1 offset $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		r.log.Warn("get all flats error: " + err.Error())
		return nil, fmt.Errorf("get all flats error:  %v", err.Error())
	}
	defer rows.Close()

	var (
		flats []domain.Flat
		flat  domain.Flat
	)
	for rows.Next() {
		err = rows.Scan(
			&flat.ID,
			&flat.HouseID,
			&flat.Price,
			&flat.Rooms,
			&flat.Status)
		if err != nil {
			r.log.Warn("get all flats error: scan flat error, " + err.Error())
			continue
		}
		flats = append(flats, flat)
	}

	return flats, err
}
