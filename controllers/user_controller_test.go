package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/ichtrojan/muzz/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserControllerTestSuite struct {
	suite.Suite

	transactionDBContainer testcontainers.Container
}

func TestUserController(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

func (p *UserControllerTestSuite) SetupSuite() {

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

func (p *UserControllerTestSuite) TearDownSuite() {
	err := p.transactionDBContainer.Terminate(context.Background())
	require.NoError(p.T(), err)
}

func (p *UserControllerTestSuite) TestCreateUser() {

	req, err := http.NewRequest("POST", "/user/create", nil)
	p.Require().NoError(err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUser)
	handler.ServeHTTP(rr, req)

	p.Require().Equal(http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	p.Require().NoError(err)

	p.Require().NotEmpty(resp["user"])
}

func (p *UserControllerTestSuite) TestLogin() {
	req, err := http.NewRequest("POST", "/user/create", nil)
	p.Require().NoError(err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUser)
	handler.ServeHTTP(rr, req)

	p.Require().Equal(http.StatusOK, rr.Code)

	var resp CreateUserResponse
	err = json.NewDecoder(rr.Body).Decode(&resp)
	p.Require().NoError(err)

	var b = new(bytes.Buffer)

	err = json.NewEncoder(b).Encode(&resp.User)
	p.Require().NoError(err)

	req, err = http.NewRequest("POST", "/user/login", b)
	p.Require().NoError(err)

	newRecorder := httptest.NewRecorder()

	handler = http.HandlerFunc(Login)
	handler.ServeHTTP(newRecorder, req)

	p.Require().Equal(http.StatusOK, newRecorder.Code)

	io.Copy(io.Discard, newRecorder.Body)

	tt := []struct {
		name        string
		requestBody struct {
			Email    string `json:"email,omitempty"`
			Password string `json:"password,omitempty"`
		} `json:"request_body,omitempty"`
		hasError           bool
		expectedStatusCode int
	}{
		{
			name: "password not provided",
			requestBody: struct {
				Email    string "json:\"email,omitempty\""
				Password string "json:\"password,omitempty\""
			}{
				Email: "trojan@trojan.com",
			},
			hasError:           true,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "email not provided",
			requestBody: struct {
				Email    string "json:\"email,omitempty\""
				Password string "json:\"password,omitempty\""
			}{
				Password: "password",
			},
			hasError:           true,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "user not found",
			requestBody: struct {
				Email    string "json:\"email,omitempty\""
				Password string "json:\"password,omitempty\""
			}{
				Email:    "trojan@trojan.com",
				Password: "password",
			},
			hasError:           true,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "password does not match",
			requestBody: struct {
				Email    string "json:\"email,omitempty\""
				Password string "json:\"password,omitempty\""
			}{
				// use the created user email so it is valid
				Email:    resp.User.Email,
				Password: "password",
			},
			hasError:           true,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "user can log in successfully",
			requestBody: struct {
				Email    string "json:\"email,omitempty\""
				Password string "json:\"password,omitempty\""
			}{
				// use the created user email and password so it is valid
				Email:    resp.User.Email,
				Password: resp.User.Password,
			},
			hasError:           false,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, v := range tt {
		p.T().Run(v.name, func(t *testing.T) {

			var b = new(bytes.Buffer)

			err = json.NewEncoder(b).Encode(&v.requestBody)
			require.NoError(t, err)

			req, err = http.NewRequest("POST", "/user/login", b)
			require.NoError(t, err)

			newRecorder := httptest.NewRecorder()

			handler = http.HandlerFunc(Login)
			handler.ServeHTTP(newRecorder, req)

			require.Equal(t, v.expectedStatusCode, newRecorder.Code)

			io.Copy(io.Discard, newRecorder.Body)
		})
	}
}
