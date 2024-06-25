package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ichtrojan/muzz/database"
	"github.com/ichtrojan/muzz/middleware"
	"github.com/ichtrojan/muzz/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type SwipeControllerTestSuite struct {
	suite.Suite
	transactionDBContainer testcontainers.Container
}

func TestSwipeController(t *testing.T) {
	suite.Run(t, new(SwipeControllerTestSuite))
}

func (p *SwipeControllerTestSuite) SetupSuite() {
	mysqlContainer, err := mysql.RunContainer(context.Background(),
		testcontainers.WithImage("mysql:8.0.36"),
		mysql.WithDatabase("foo"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
		mysql.WithScripts(filepath.Join("testdata", "schema.sql")),
	)

	p.Require().NoError(err)

	p.transactionDBContainer = mysqlContainer

	port, err := p.transactionDBContainer.MappedPort(context.Background(), "3306")
	p.Require().NoError(err)

	err = database.ConnectMySQL("root", "password", "localhost", port.Port(), "foo")
	p.Require().NoError(err)
}

func (p *SwipeControllerTestSuite) TearDownSuite() {
	err := p.transactionDBContainer.Terminate(context.Background())
	require.NoError(p.T(), err)
}

func (p *SwipeControllerTestSuite) TestSwipe() {
	req, err := http.NewRequest("POST", "/user/create", nil)
	p.Require().NoError(err)

	var users = make([]models.User, 0)

	// Create two users
	// Swipe against each other

	for range []int{0, 1} {

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateUser)
		handler.ServeHTTP(rr, req)

		p.Require().Equal(http.StatusOK, rr.Code)

		var resp CreateUserResponse
		err = json.NewDecoder(rr.Body).Decode(&resp)
		p.Require().NoError(err)

		var swipedUser models.User
		err = database.Connection.Where("email = ?", resp.User.Email).First(&swipedUser).Error

		p.Require().NoError(err)

		users = append(users, models.User{
			Id:       swipedUser.Id,
			Email:    resp.User.Email,
			Password: resp.User.Password,
		})
	}

	// sanity checks
	p.Require().Len(users, 2)

	// for each user, swipe against each other
	// Login first
	// then swipe
	for idx, user := range users {
		var b = new(bytes.Buffer)

		err = json.NewEncoder(b).Encode(&user)
		p.Require().NoError(err)

		req, err = http.NewRequest("POST", "/user/login", b)
		p.Require().NoError(err)

		newRecorder := httptest.NewRecorder()

		handler := http.HandlerFunc(Login)
		handler.ServeHTTP(newRecorder, req)

		p.Require().Equal(http.StatusOK, newRecorder.Code)

		var resp struct {
			Token string `json:"token"`
		}

		err = json.NewDecoder(newRecorder.Body).Decode(&resp)
		p.Require().NoError(err)
		p.Require().NotEmpty(resp.Token)

		var idxToUse = 0
		if idx == 0 {
			idxToUse = 1
		}

		var reqbody = struct {
			UserId     string
			Preference string
		}{
			UserId:     users[idxToUse].Id,
			Preference: "yes",
		}

		b.Reset()

		err = json.NewEncoder(b).Encode(&reqbody)
		p.Require().NoError(err)

		req, err = http.NewRequest("POST", "/user/swipe", b)
		p.Require().NoError(err)

		req.Header.Add("Authorization", "Bearer "+resp.Token)

		// Swipe now
		// Make sure to run the authentication middleware
		swipeHandler := middleware.AuthenticateUser(http.HandlerFunc(Swipe))
		swipeRecorder := httptest.NewRecorder()
		swipeHandler.ServeHTTP(swipeRecorder, req)

		p.Require().Equal(http.StatusOK, newRecorder.Code)
		io.Copy(os.Stdout, swipeRecorder.Body)
	}
}
