package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Edgar200021/netowork-server-go/tests/testapp"
	"github.com/alecthomas/assert/v2"
	"github.com/go-faker/faker/v4"
)

func TestResetPassword(t *testing.T) {
	setup := func(t *testing.T) (*testapp.TestApp, testapp.SignupData) {
		t.Helper()

		app := testapp.New(t)
		signupData := testapp.GenerateFakeData[testapp.SignupData](t)

		app.CreateAndVerify(t, signupData)

		response, err := app.ForgotPassword(
			map[string]string{
				"email": signupData.Email,
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode())

		return app, signupData
	}
	t.Parallel()

	t.Run(
		"Returns 200 when request is valid",
		func(t *testing.T) {
			t.Parallel()
			app, data := setup(t)

			token, err := app.GetResetPasswordToken()
			assert.NoError(t, err)
			assert.NotEqual(t, "", token)

			response, err := app.ResetPassword(
				map[string]string{
					"email":    data.Email,
					"token":    token,
					"password": strings.Repeat("q", testapp.PasswordMinLength),
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode())
		},
	)

	t.Run(
		"Should apply changes into database when request is valid",
		func(t *testing.T) {
			t.Parallel()
			app, data := setup(t)

			token, err := app.GetResetPasswordToken()
			assert.NoError(t, err)
			assert.NotEqual(t, "", token)

			userBefore, err := app.GetUser(data.Email)
			assert.NoError(t, err)

			response, err := app.ResetPassword(
				map[string]string{
					"email":    data.Email,
					"token":    token,
					"password": strings.Repeat("q", testapp.PasswordMinLength),
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode())

			userAfter, err := app.GetUser(data.Email)
			assert.NoError(t, err)

			assert.True(t, userAfter.IsVerified)
			assert.NotEqual(t, userBefore.Password, userAfter.Password)
		},
	)

	t.Run(
		"Returns 400 when request is invalid",
		func(t *testing.T) {
			t.Parallel()
			app, _ := setup(t)

			testCases := []struct {
				name               string
				data               map[string]string
				expectedErrorField string
			}{
				{
					name: "empty email",
					data: map[string]string{
						"email":    "",
						"token":    "token",
						"password": strings.Repeat("q", testapp.PasswordMinLength),
					},
					expectedErrorField: "email",
				},
				{
					name: "invalid email",
					data: map[string]string{
						"email":    "invalid email",
						"token":    "token",
						"password": strings.Repeat("q", testapp.PasswordMinLength),
					},
					expectedErrorField: "email",
				},
				{
					name: "empty token",
					data: map[string]string{
						"email":    faker.Email(),
						"token":    "",
						"password": strings.Repeat("q", testapp.PasswordMinLength),
					},
					expectedErrorField: "token",
				},
				{
					name: "empty password",
					data: map[string]string{
						"email":    faker.Email(),
						"token":    "token",
						"password": "",
					},
					expectedErrorField: "password",
				},
				{
					name: "password too short",
					data: map[string]string{
						"email":    faker.Email(),
						"token":    "token",
						"password": strings.Repeat("q", testapp.PasswordMinLength-1),
					},
					expectedErrorField: "password",
				},
				{
					name: "password too long",
					data: map[string]string{
						"email":    faker.Email(),
						"token":    "token",
						"password": strings.Repeat("q", testapp.PasswordMaxLength+1),
					},
					expectedErrorField: "password",
				},
			}

			for _, tc := range testCases {
				tc := tc

				t.Run(
					tc.name,
					func(t *testing.T) {
						t.Parallel()

						response, err := app.ResetPassword(tc.data)

						assert.NoError(t, err)
						assert.Equal(t, http.StatusBadRequest, response.StatusCode())
						app.AssertValidationErrors(t, response, tc.expectedErrorField)
					},
				)
			}
		},
	)

	t.Run(
		"Returns 400 when user is banned",
		func(t *testing.T) {
			t.Parallel()
			app, data := setup(t)

			err := app.BanUser(data.Email)
			assert.NoError(t, err)

			token, err := app.GetResetPasswordToken()
			assert.NoError(t, err)
			assert.NotEqual(t, "", token)

			response, err := app.ResetPassword(
				map[string]string{
					"email":    data.Email,
					"token":    token,
					"password": strings.Repeat("q", testapp.PasswordMaxLength),
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		},
	)

	t.Run(
		"Returns 400 when email is different",
		func(t *testing.T) {
			t.Parallel()
			app, _ := setup(t)

			token, err := app.GetResetPasswordToken()
			assert.NoError(t, err)
			assert.NotEqual(t, "", token)

			response, err := app.ResetPassword(
				map[string]string{
					"email":    faker.Email(),
					"token":    token,
					"password": strings.Repeat("q", testapp.PasswordMaxLength),
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		},
	)
}
