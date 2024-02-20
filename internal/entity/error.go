package entity

import (
	"fmt"
	"time"
)

type ErrRoomNotAvailable struct {
	UnavailableDates []time.Time
}

func (e ErrRoomNotAvailable) Error() string {
	return fmt.Sprintf("room not available: %v", e.UnavailableDates)
}

// To use in handler layer
//func (e ErrRoomNotAvailable) HTTPError() HTTPError {
//
//}
