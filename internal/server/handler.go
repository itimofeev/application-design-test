package server

import (
	"net/http"

	"applicationDesignTest/internal/app/booking"
	"applicationDesignTest/internal/entity"
)

type handler struct {
	bookingApp *booking.App
}

func (h *handler) createOrder(w http.ResponseWriter, r *http.Request) {
	// todo parse OrderRequest
	orderRequest := entity.OrderRequest{}

	order, err := h.bookingApp.CreateOrder(r.Context(), orderRequest)
	if err != nil {
		// todo call error handler
		return
	}

	_ = order
	// todo write order to response
}
