package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"wb_order_service/internal/cache"
	"wb_order_service/internal/model"

	"github.com/go-chi/chi/v5"
)

type fakeStore struct {
	orders map[string]*model.Order
}

func (f *fakeStore) GetOrderByID(_ context.Context, id string) (*model.Order, error) {
	if o, ok := f.orders[id]; ok {
		return o, nil
	}
	return nil, errNotFound
}
func (f *fakeStore) GetAllOrders(context.Context) ([]*model.Order, error) { return nil, nil }
func (f *fakeStore) SaveOrder(context.Context, *model.Order) error        { return nil }

var errNotFound = fmt.Errorf("order not found")

func TestGetOrderByID_OK(t *testing.T) {
	c := cache.NewCache()
	o := &model.Order{OrderUID: "abc123", TrackNumber: "TN", Entry: "E", CustomerID: "C", DeliveryService: "D", ShardKey: "S", DateCreated: model.Order{}.DateCreated, OofShard: "1"}
	c.Set(o)
	h := NewHandler(c, &fakeStore{orders: map[string]*model.Order{"abc123": o}})

	req := httptest.NewRequest(http.MethodGet, "/order/abc123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("order_uid", "abc123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()
	h.GetOrderByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestGetOrderByID_BadRequest(t *testing.T) {
	h := NewHandler(cache.NewCache(), &fakeStore{orders: map[string]*model.Order{}})
	req := httptest.NewRequest(http.MethodGet, "/order/", nil)
	rr := httptest.NewRecorder()
	h.GetOrderByID(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}
