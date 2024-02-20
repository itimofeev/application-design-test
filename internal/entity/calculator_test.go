package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMakeOrder(t *testing.T) {
	tests := []struct {
		name string

		availability []RoomAvailability
		request      OrderRequest

		expectedOrder    Order
		expectedNewAvail []RoomAvailability
		expectedError    error
	}{
		{
			name: "Valid order",
			availability: []RoomAvailability{
				{HotelID: "hotel1", RoomID: "room1", Date: date(2024, 1, 1), Quota: 5},
				{HotelID: "hotel2", RoomID: "room2", Date: date(2024, 1, 2), Quota: 3},
			},
			request: OrderRequest{
				HotelID:   "hotel1",
				RoomID:    "room1",
				UserEmail: "user1",
				From:      date(2024, 1, 1),
				To:        date(2024, 1, 2),
				Count:     2,
			},
			expectedOrder: Order{
				HotelID:   "hotel1",
				RoomID:    "room1",
				UserEmail: "user1",
				From:      date(2024, 1, 1),
				To:        date(2024, 1, 2),
			},
			expectedNewAvail: []RoomAvailability{
				{HotelID: "hotel1", RoomID: "room1", Date: date(2024, 1, 1), Quota: 3},
				{HotelID: "hotel2", RoomID: "room2", Date: date(2024, 1, 2), Quota: 1},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := OrderCalculator{Availability: tt.availability, Request: tt.request}
			order, newAvailability, err := calculator.MakeOrder()

			if tt.expectedError != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedError.Error())
				return
			}
			require.NoError(t, err)

			require.Equal(t, tt.expectedOrder, order)
			require.Equal(t, tt.expectedNewAvail, newAvailability)
		})
	}
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
