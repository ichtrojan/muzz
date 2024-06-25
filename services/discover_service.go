package services

import (
	"github.com/ichtrojan/muzz/database"
	"github.com/ichtrojan/muzz/helpers"
	"github.com/ichtrojan/muzz/middleware"
	"github.com/ichtrojan/muzz/models"
	"net/http"
	"strconv"
)

func Discover(r *http.Request) ([]models.User, helpers.PaginateData, error, int) {
	user := middleware.GetUser(r.Context())

	var users []models.User

	var totalRecords int64

	query := database.Connection.Scopes(helpers.Paginate(r)).Where("id != ?", user.Id)

	ageParam := r.URL.Query().Get("age")

	gender := r.URL.Query().Get("gender")

	if ageParam != "" {
		age, err := strconv.Atoi(ageParam)
		if err != nil {
			return []models.User{}, helpers.PaginateData{}, helpers.ServerError(), http.StatusInternalServerError
		}

		query = query.Where("age = ?", age)
	}

	if gender != "" {
		query = query.Where("gender = ?", gender)
	}

	var swipedIds []string

	database.Connection.Model(&models.Swipe{}).Where("user_id = ?", user.Id).Select("swiped_on").Find(&swipedIds)

	if len(swipedIds) > 0 {
		query = query.Where("id NOT IN ?", swipedIds)
	}

	_ = query.Find(&users).Count(&totalRecords)

	pagination := helpers.GetPagination(r, int(totalRecords))

	return users, pagination, nil, http.StatusOK
}
