package domain

import (
	"errors"
	"time"
)

var (
	ErrHouseBadID      = errors.New("bad house id")
	ErrHouseBadAddress = errors.New("bad house address")
	ErrHouseBadYear    = errors.New("bad house year")
)

type House struct {
	HouseID   int
	Year      int
	Address   string
	Developer string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type HouseCreateRequest struct {
	Address   string `json:"address"`
	Year      int    `json:"year"`
	Developer string `json:"developer"`
}

func (h *HouseCreateRequest) Validate() error {
	if h.Address == "" {
		return ErrHouseBadAddress
	}
	if h.Year < 0 {
		return ErrHouseBadYear
	}
	return nil
}

type HouseCreateResponse struct {
	HomeID    int       `json:"id"`
	Address   string    `json:"address"`
	Year      int       `json:"year"`
	Developer string    `json:"developer"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
}

type HouseFlatsRequest struct {
	ID int `json:"id"`
}

type HouseFlatsResponse struct {
	Flats []HouseFlatResponse `json:"flats"`
}

type HouseFlatResponse struct {
	ID      int    `json:"id"`
	HouseID int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"`
}
