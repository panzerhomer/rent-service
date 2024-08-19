package domain

import "github.com/google/uuid"

const (
	StatusCreated      = "created"
	StatusApproved     = "approved"
	StatusDeclined     = "declined"
	StatusOnModeration = "on_moderation"
)

type Flat struct {
	ID          int
	HouseID     int
	UserID      uuid.UUID
	Price       int
	Rooms       int
	Status      string
	ModeratorID int
}

type FlatCreateRequest struct {
	HouseID int `json:"house_id"`
	Price   int `json:"price"`
	Rooms   int `json:"rooms"`
}

type FlatCreateResponse struct {
	ID      int    `json:"id"`
	HouseID int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"`
}

type FlatUpdateRequest struct {
	ID      int    `json:"id"`
	HouseID int    `json:"house_id"`
	Status  string `json:"status,omitempty"`
}
