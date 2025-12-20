package dtos

type CreateOrderResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserID   int    `json:"user_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type GetOrderResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserID   int    `json:"user_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type GetAllOrdersResponse struct {
	Orders []GetOrderResponse `json:"orders"`
}

type UpdateOrderResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserID   int    `json:"user_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type DeleteOrderResponse struct {
	Success bool `json:"success"`
}
