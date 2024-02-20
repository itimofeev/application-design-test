package booking

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"applicationDesignTest/internal/entity"
)

type repository interface {
	GetAvailability(ctx context.Context, filter entity.RoomSearchFilter) ([]entity.RoomAvailability, error)
	UpdateAvailability(ctx context.Context, availability []entity.RoomAvailability) error
	AddOrder(ctx context.Context, order entity.Order) error
	DoInTx(ctx context.Context, f func(ctx context.Context) error) error
}

type Config struct {
	Repository repository
}

type App struct {
	repository repository
}

func New(cfg Config) (*App, error) {
	err := validator.New().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		repository: cfg.Repository,
	}, nil
}

func (a *App) CreateOrder(ctx context.Context, req entity.OrderRequest) (entity.Order, error) {
	// todo add logs
	var order entity.Order
	err := a.repository.DoInTx(ctx, func(ctx context.Context) error {
		availability, err := a.repository.GetAvailability(ctx, req.ToFilter())
		if err != nil {
			return fmt.Errorf("failed to get availability: %w", err)
		}

		calculator := entity.OrderCalculator{
			Availability: availability,
			Request:      req,
		}

		var newAvailability []entity.RoomAvailability
		order, newAvailability, err = calculator.MakeOrder()

		if err != nil {
			return fmt.Errorf("failed to check availability: %w", err)
		}

		err = a.repository.UpdateAvailability(ctx, newAvailability)
		if err != nil {
			return fmt.Errorf("failed to update availability: %w", err)
		}

		err = a.repository.AddOrder(ctx, order)
		if err != nil {
			return fmt.Errorf("failed to add order: %w", err)
		}

		return nil
	})

	if err != nil {
		return entity.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}
