package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"
	. "github.com/Edgar200021/netowork-server-go/.gen/netowork/public/table"
	"github.com/Edgar200021/netowork-server-go/tests/testapp"
	"github.com/alecthomas/assert/v2"
	. "github.com/go-jet/jet/v2/postgres"
)

func TestSignup(t *testing.T) {
	setup := func(t *testing.T) (*testapp.TestApp, testapp.SignupData) {
		t.Helper()
		return testapp.New(t), testapp.GenerateFakeData[testapp.SignupData](t)
	}
	t.Parallel()

	t.Run(
		"Returns 200 status code when request is successful", func(t *testing.T) {
			t.Parallel()
			app, data := setup(t)

			response, err := app.SignUp(data)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, response.StatusCode())
		},
	)

	t.Run(
		"Should be saved into database when request is successful",
		func(t *testing.T) {
			t.Parallel()
			app, data := setup(t)

			response, err := app.SignUp(data)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, response.StatusCode())

			var dbUser model.Users

			stmt := SELECT(
				Users.ID,
				Users.IsVerified,
				Users.Password,
			).
				FROM(Users).
				WHERE(Users.Email.EQ(Text(data.Email)))

			err = stmt.Query(app.Db, &dbUser)

			assert.NoError(t, err)
			assert.NotEqual(t, "", dbUser.ID.String())
			assert.NotEqual(t, data.Password, dbUser.Password)
			assert.False(t, dbUser.IsVerified)
		},
	)

	t.Run(
		"Returns 400 when request is invalid",
		func(t *testing.T) {
			t.Parallel()
			app, data := setup(t)

			testCases := []struct {
				name               string
				data               testapp.SignupData
				expectedErrorField string
			}{
				{
					name: "empty email",
					data: func() testapp.SignupData {
						d := data
						d.Email = ""
						return d
					}(),
					expectedErrorField: "email",
				},
				{
					name: "invalid email",
					data: func() testapp.SignupData {
						d := data
						d.Email = "invalid"
						return d
					}(),
					expectedErrorField: "email",
				},
				{
					name: "empty password",
					data: func() testapp.SignupData {
						d := data
						d.Password = ""
						return d
					}(),
					expectedErrorField: "password",
				},
				{
					name: "password too short",
					data: func() testapp.SignupData {
						d := data
						d.Password = strings.Repeat("q", testapp.PasswordMinLength-1)
						return d
					}(),
					expectedErrorField: "password",
				},
				{
					name: "password too long",
					data: func() testapp.SignupData {
						d := data
						d.Password = strings.Repeat("q", testapp.PasswordMaxLength+1)
						return d
					}(),
					expectedErrorField: "password",
				},
				{
					name: "empty firstName",
					data: func() testapp.SignupData {
						d := data
						d.FirstName = ""
						return d
					}(),
					expectedErrorField: "firstName",
				},
				{
					name: "firstName too short",
					data: func() testapp.SignupData {
						d := data
						d.FirstName = strings.Repeat("q", testapp.FirstNameMinLength-1)
						return d
					}(),
					expectedErrorField: "firstName",
				},
				{
					name: "firstName too long",
					data: func() testapp.SignupData {
						d := data
						d.FirstName = strings.Repeat("q", testapp.FirstNameMaxLength+1)
						return d
					}(),
					expectedErrorField: "firstName",
				},
				{
					name: "empty lastName",
					data: func() testapp.SignupData {
						d := data
						d.LastName = ""
						return d
					}(),
					expectedErrorField: "lastName",
				},
				{
					name: "lastName too short",
					data: func() testapp.SignupData {
						d := data
						d.LastName = strings.Repeat("q", testapp.LastNameMinLength-1)
						return d
					}(),
					expectedErrorField: "lastName",
				},
				{
					name: "lastName too long",
					data: func() testapp.SignupData {
						d := data
						d.LastName = strings.Repeat("q", testapp.LastNameMaxLength+1)
						return d
					}(),
					expectedErrorField: "lastName",
				},
				{
					name: "empty role",
					data: func() testapp.SignupData {
						d := data
						d.Role = ""
						return d
					}(),
					expectedErrorField: "role",
				},
				{
					name: "invalid role",
					data: func() testapp.SignupData {
						d := data
						d.Role = "invalidrole"
						return d
					}(),
					expectedErrorField: "role",
				},
				{
					name: "admin role",
					data: func() testapp.SignupData {
						d := data
						d.Role = "admin"
						return d
					}(),
					expectedErrorField: "role",
				},
			}

			for _, test := range testCases {
				t.Run(
					test.name, func(t *testing.T) {
						t.Parallel()
						response, err := app.SignUp(test.data)
						assert.NoError(t, err)

						assert.Equal(t, http.StatusBadRequest, response.StatusCode())
						app.AssertValidationErrors(t, response, test.expectedErrorField)
					},
				)
			}
		},
	)

	t.Run(
		"Returns 400 when user already exists",
		func(t *testing.T) {
			t.Parallel()
			app, data := setup(t)

			firstResponse, err := app.SignUp(data)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, firstResponse.StatusCode())

			secondResponse, err := app.SignUp(data)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, secondResponse.StatusCode())
		},
	)

}
