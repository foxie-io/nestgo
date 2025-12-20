package dtos

type CreateOrderRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserID   int    `json:"user_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type UpdateOrderRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserID   int    `json:"user_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type GetOrderRequest struct {
	ID int `json:"id" uri:"id"`
}

type DeleteOrderRequest struct {
	ID int `json:"id" uri:"id"`
}

type ListOrdersRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
