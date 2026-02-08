package adapters

import (
	"context"
	"errors"
	"sync"

	"github.com/MousaZa/logistics-management/internal/orders/domain/orders"
)

type OrderMemoryRepository struct {
	orders []*orders.Order
	mutex  sync.RWMutex
}

func NewOrderMemoryRepository() *OrderMemoryRepository {
	return &OrderMemoryRepository{}
}

func (mr *OrderMemoryRepository) AddOrder(_ context.Context, o *orders.Order) error {
	mr.mutex.Lock()
	defer mr.mutex.Unlock()
	mr.orders = append(mr.orders, o)
	return nil
}

func (mr *OrderMemoryRepository) GetOrder(_ context.Context, orderUUID string) (*orders.Order, error) {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()
	for _, o := range mr.orders {
		if o.OrderUUID == orderUUID {
			return o, nil
		}
	}
	return nil, errors.New("order not found")
}

func (mr *OrderMemoryRepository) GetAllOrders(_ context.Context) ([]*orders.Order, error) {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()
	return mr.orders, nil
}

func (mr *OrderMemoryRepository) UpdateOrder(ctx context.Context, orderUUID string, updateFunc func(ctx context.Context, o *orders.Order) (*orders.Order, error)) error {
	mr.mutex.Lock()
	defer mr.mutex.Unlock()
	for i, o := range mr.orders {
		if o.OrderUUID == orderUUID {
			no, err := updateFunc(ctx, o)
			if err != nil {
				return err
			}
			mr.orders[i] = no
			return nil
		}
	}
	return errors.New("order not found")
}
