package dtos

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetUserRequest struct {
	ID int `json:"id" uri:"id"`
}

type DeleteUserRequest struct {
	ID int `json:"id" uri:"id"`
}

type ListUsersRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
