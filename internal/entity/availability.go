package entity

import (
	"errors"
	"time"
)

type RoomAvailability struct {
	HotelID string
	RoomID  string
	Date    time.Time
	Quota   uint64
}

func (a RoomAvailability) ApplyOrder(req OrderRequest) (RoomAvailability, error) {
	if a.Quota < req.Count {
		return RoomAvailability{}, errors.New("room is not available")
	}

	return RoomAvailability{
		HotelID: a.HotelID,
		RoomID:  a.RoomID,
		Date:    a.Date,
		Quota:   a.Quota - req.Count,
	}, nil
}
