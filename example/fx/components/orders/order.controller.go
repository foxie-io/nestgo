package orders

import (
	"context"
	"example/fx/adapter"
	"example/fx/components/orders/dtos"

	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
)

type OrderController struct {
	ng.DefaultControllerInitializer
	order_s *OrderService
}

func NewOrderController(order_s *OrderService) *OrderController {
	return &OrderController{
		order_s: order_s,
	}
}

func (con *OrderController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/orders"),
	)
}

func (con *OrderController) Create() ng.Route {
	return ng.NewRoute(http.MethodPost, "/",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				body dtos.CreateOrderRequest
			)

			return ng.Handle(
				adapter.BindBody(&body),
				func(ctx context.Context) error {
					resp, err := con.order_s.CreateOrder(body)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

func (con *OrderController) Get() ng.Route {
	return ng.NewRoute(http.MethodGet, "/:id",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				param dtos.GetOrderRequest
			)
			return ng.Handle(
				adapter.BindParams(&param),
				func(ctx context.Context) error {
					resp, err := con.order_s.GetOrder(param.ID)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

func (con *OrderController) GetAll() ng.Route {
	return ng.NewRoute(http.MethodGet, "/",
		ng.WithScopeHandler(func() ng.Handler {
			return ng.Handle(
				func(ctx context.Context) error {
					resp := con.order_s.GetAllOrders()
					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

func (con *OrderController) Update() ng.Route {
	return ng.NewRoute(http.MethodPut, "/:id",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				param dtos.DeleteOrderRequest
				body  dtos.UpdateOrderRequest
			)

			return ng.Handle(
				adapter.BindBody(&body),
				adapter.BindParams(&param),
				func(ctx context.Context) error {
					resp, err := con.order_s.UpdateOrder(param.ID, &body)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

func (con *OrderController) Delete() ng.Route {
	return ng.NewRoute(http.MethodDelete, "/:id",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				param dtos.DeleteOrderRequest
			)
			return ng.Handle(
				adapter.BindParams(&param),
				func(ctx context.Context) error {
					resp := con.order_s.DeleteOrder(&param)
					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}
