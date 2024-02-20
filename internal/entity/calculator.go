package entity

import (
	"time"
)

type OrderCalculator struct {
	Availability []RoomAvailability
	Request      OrderRequest
}

func (c OrderCalculator) MakeOrder() (Order, []RoomAvailability, error) {
	unavailableDays := make([]time.Time, 0)
	newAvailability := make([]RoomAvailability, 0)
	for _, availability := range c.Availability {
		newRoomAvailability, err := availability.ApplyOrder(c.Request)
		if err != nil {
			unavailableDays = append(unavailableDays, availability.Date)
			continue
		}

		newAvailability = append(newAvailability, newRoomAvailability)
	}

	if len(unavailableDays) != 0 {
		return Order{}, newAvailability, ErrRoomNotAvailable{UnavailableDates: unavailableDays}
	}

	order := Order{
		HotelID:   c.Request.HotelID,
		RoomID:    c.Request.RoomID,
		UserEmail: c.Request.UserEmail,
		From:      c.Request.From,
		To:        c.Request.To,
		Price:     0, // calculate price taking into account bonuses and price info
	}
	return order, newAvailability, nil
}
