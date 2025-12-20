package orders

import (
	"example/fx/components/orders/dtos"
	"example/fx/models"

	nghttp "github.com/foxie-io/ng/http"
)

type OrderService struct {
	orders    map[int]*models.Order
	orderList []*models.Order
}

func NewOrderService() *OrderService {
	return &OrderService{
		orders:    make(map[int]*models.Order),
		orderList: []*models.Order{},
	}
}

func (s *OrderService) CreateOrder(req dtos.CreateOrderRequest) (*dtos.CreateOrderResponse, error) {
	id := len(s.orders) + 1
	order := &models.Order{
		ID:       id,
		UserID:   req.UserID,
		Product:  req.Product,
		Quantity: req.Quantity,
	}
	s.orders[id] = order
	s.orderList = append(s.orderList, order)
	return &dtos.CreateOrderResponse{
		ID:       id,
		UserID:   order.UserID,
		Product:  order.Product,
		Quantity: order.Quantity,
	}, nil
}

func (s *OrderService) GetOrder(id int) (*dtos.GetOrderResponse, error) {
	order, exists := s.orders[id]
	if !exists {
		return nil, nghttp.NewErrNotFound()
	}
	return &dtos.GetOrderResponse{
		ID:       order.ID,
		UserID:   order.UserID,
		Product:  order.Product,
		Quantity: order.Quantity,
	}, nil
}

func (s *OrderService) GetAllOrders() *dtos.GetAllOrdersResponse {
	var orders []dtos.GetOrderResponse
	for _, order := range s.orderList {
		orders = append(orders, dtos.GetOrderResponse{
			ID:       order.ID,
			UserID:   order.UserID,
			Product:  order.Product,
			Quantity: order.Quantity,
		})
	}
	return &dtos.GetAllOrdersResponse{Orders: orders}
}

func (s *OrderService) UpdateOrder(id int, req *dtos.UpdateOrderRequest) (*dtos.UpdateOrderResponse, error) {
	if _, exists := s.orders[id]; !exists {
		return nil, nghttp.NewErrNotFound()
	}

	updatedOrder := &models.Order{
		ID:       id,
		UserID:   req.UserID,
		Product:  req.Product,
		Quantity: req.Quantity,
	}
	s.orders[id] = updatedOrder

	for i, order := range s.orderList {
		if order.ID == id {
			s.orderList[i] = updatedOrder
			break
		}
	}

	return &dtos.UpdateOrderResponse{
		ID:       updatedOrder.ID,
		UserID:   updatedOrder.UserID,
		Product:  updatedOrder.Product,
		Quantity: updatedOrder.Quantity,
	}, nil
}

func (s *OrderService) DeleteOrder(params *dtos.DeleteOrderRequest) *dtos.DeleteOrderResponse {
	if _, exists := s.orders[params.ID]; !exists {
		return &dtos.DeleteOrderResponse{Success: false}
	}
	delete(s.orders, params.ID)

	for i, order := range s.orderList {
		if order.ID == params.ID {
			s.orderList = append(s.orderList[:i], s.orderList[i+1:]...)
			break
		}
	}
	return &dtos.DeleteOrderResponse{Success: true}
}
