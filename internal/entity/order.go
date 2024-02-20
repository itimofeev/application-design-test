package entity

import (
	"time"
)

type Order struct {
	HotelID   string
	RoomID    string
	UserEmail string
	From      time.Time
	To        time.Time

	Price int
}

type OrderRequest struct {
	HotelID   string
	RoomID    string
	UserEmail string
	From      time.Time
	To        time.Time
	Count     uint64
}

func (r OrderRequest) ToFilter() RoomSearchFilter {
	return RoomSearchFilter{
		HotelID: r.HotelID,
		RoomID:  r.RoomID,
		From:    r.From,
		To:      r.To,
	}
}

type RoomSearchFilter struct {
	HotelID string
	RoomID  string
	From    time.Time
	To      time.Time
}
