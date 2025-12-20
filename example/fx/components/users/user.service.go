package users

import (
	"example/fx/components/users/dtos"
	"example/fx/models"

	nghttp "github.com/foxie-io/ng/http"
)

type UserService struct {
	users    map[int]*models.User
	userList []*models.User
}

func NewUserService() *UserService {
	return &UserService{
		users:    make(map[int]*models.User),
		userList: []*models.User{},
	}
}

func (s *UserService) CreateUser(req dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	id := len(s.users) + 1
	user := &models.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	s.users[id] = user
	s.userList = append(s.userList, user)
	return &dtos.CreateUserResponse{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) GetUser(id int) (*dtos.GetUserResponse, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, nghttp.NewErrNotFound()
	}
	return &dtos.GetUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) GetAllUsers() *dtos.GetAllUsersResponse {
	var users []dtos.GetUserResponse
	for _, user := range s.userList {
		users = append(users, dtos.GetUserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return &dtos.GetAllUsersResponse{Users: users}
}

func (s *UserService) UpdateUser(id int, req *dtos.UpdateUserRequest) (*dtos.UpdateUserResponse, error) {
	if _, exists := s.users[id]; !exists {
		return nil, nghttp.NewErrNotFound()
	}

	updatedUser := &models.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	s.users[id] = updatedUser

	for i, user := range s.userList {
		if user.ID == id {
			s.userList[i] = updatedUser
			break
		}
	}

	return &dtos.UpdateUserResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	}, nil
}

func (s *UserService) DeleteUser(params *dtos.DeleteUserRequest) *dtos.DeleteUserResponse {
	if _, exists := s.users[params.ID]; !exists {
		return &dtos.DeleteUserResponse{Success: false}
	}
	delete(s.users, params.ID)

	for i, user := range s.userList {
		if user.ID == params.ID {
			s.userList = append(s.userList[:i], s.userList[i+1:]...)
			break
		}
	}
	return &dtos.DeleteUserResponse{Success: true}
}
