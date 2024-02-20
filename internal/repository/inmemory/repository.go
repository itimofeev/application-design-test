package inmemory

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"

	"applicationDesignTest/internal/entity"
)

type Config struct {
}

type Repository struct {
	orders       []entity.Order
	availability []entity.RoomAvailability
	// todo add mutex
}

func New(cfg Config) (*Repository, error) {
	err := validator.New().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return &Repository{
		availability: []entity.RoomAvailability{
			{"reddison", "lux", date(2024, 1, 1), 1},
			{"reddison", "lux", date(2024, 1, 2), 1},
			{"reddison", "lux", date(2024, 1, 3), 1},
			{"reddison", "lux", date(2024, 1, 4), 1},
			{"reddison", "lux", date(2024, 1, 5), 0},
		},
		orders: []entity.Order{},
	}, nil
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func (r *Repository) GetAvailability(ctx context.Context, filter entity.RoomSearchFilter) ([]entity.RoomAvailability, error) {
	// todo apply filters
	return r.availability, nil
}

func (r *Repository) UpdateAvailability(ctx context.Context, availability []entity.RoomAvailability) error {
	return nil
}

// AddOrder saves order in orders and events table.
// Events will be sent to kafka in separate worker implementing transactional outbox pattern.
func (r *Repository) AddOrder(ctx context.Context, order entity.Order) error {
	return nil
}

func (r *Repository) DoInTx(ctx context.Context, f func(ctx context.Context) error) error {
	// todo begin transaction and put it in context
	txCtx := context.WithValue(ctx, "tx", true)
	// todo defer rollback transaction if error

	err := f(txCtx)
	if err != nil {
		return err
	}
	return nil
}
