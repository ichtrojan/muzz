package responses

import (
	"github.com/ichtrojan/muzz/helpers"
	"github.com/ichtrojan/muzz/models"
)

type UserResponse struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Age       int    `json:"age"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type FakeUserResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

func GenerateUserResponse(user models.User) UserResponse {
	return UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Gender:    user.Gender,
		Age:       user.Age,
		CreatedAt: helpers.JSONTime{Time: user.CreatedAt}.Json(),
		UpdatedAt: helpers.JSONTime{Time: user.UpdatedAt}.Json(),
	}
}

func GenerateFakeUserResponse(user models.User) FakeUserResponse {
	return FakeUserResponse{
		Id:       user.Id,
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
		Gender:   user.Gender,
		Age:      user.Age,
	}
}

func GenerateUsersResponse(users []models.User) []UserResponse {
	var response []UserResponse

	for _, user := range users {
		response = append(response, GenerateUserResponse(user))
	}

	if len(response) == 0 {
		return []UserResponse{}
	}

	return response
}
