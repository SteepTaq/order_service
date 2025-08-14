package database

import (
	"context"
	"wb_order_service/internal/model"
)

type OrderStore interface {
	GetOrderByID(ctx context.Context, orderUID string) (*model.Order, error)
	GetAllOrders(ctx context.Context) ([]*model.Order, error)
	SaveOrder(ctx context.Context, order *model.Order) error
}
