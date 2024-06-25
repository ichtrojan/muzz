package responses

import (
	"context"
	"github.com/ichtrojan/muzz/helpers"
	"github.com/ichtrojan/muzz/middleware"
	"github.com/ichtrojan/muzz/models"
)

type DiscoverPaginatedResponse struct {
	Users []DiscoverResponse   `json:"users"`
	Meta  helpers.PaginateData `json:"meta"`
}

type DiscoverResponse struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	Gender         string  `json:"gender"`
	Age            int     `json:"age"`
	DistanceFromMe float64 `json:"distanceFromMe"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

func GenerateDiscoverResponse(ctx context.Context, users []models.User) []DiscoverResponse {
	var response []DiscoverResponse

	for _, user := range users {
		distance, _ := helpers.GetDistance(middleware.GetUser(ctx), user)

		response = append(response, DiscoverResponse{
			Id:             user.Id,
			Name:           user.Name,
			DistanceFromMe: distance,
			Gender:         user.Gender,
			Age:            user.Age,
			CreatedAt:      helpers.JSONTime{Time: user.CreatedAt}.Json(),
			UpdatedAt:      helpers.JSONTime{Time: user.UpdatedAt}.Json(),
		})
	}

	if len(response) == 0 {
		return []DiscoverResponse{}
	}

	return response
}
