package helpers

import (
	"github.com/ichtrojan/muzz/database"
	"github.com/ichtrojan/muzz/models"
)

type Result struct {
	Distance float64
}

func GetDistance(firstUser, secondUser models.User) (float64, error) {
	var result Result

	query := `SELECT ST_Distance_Sphere(point(?, ?), point(?, ?)) AS distance`

	err := database.Connection.Raw(query, firstUser.Longitude, firstUser.Latitude, secondUser.Longitude, secondUser.Latitude).Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.Distance, nil
}
