package auth

import (
	"net/http"
	"testing"

	"github.com/Edgar200021/netowork-server-go/tests/testapp"
	"github.com/alecthomas/assert/v2"
)

func TestForgotPassword(t *testing.T) {
	setup := func(t *testing.T) (*testapp.TestApp, testapp.SignupData) {
		t.Helper()
		return testapp.New(t), testapp.GenerateFakeData[testapp.SignupData](t)
	}
	t.Parallel()
	
	t.Run(
		"Returns 200 when request is valid",
		func(t *testing.T) {
			t.Parallel()
			app, data := setup(t)
			app.CreateAndVerify(t, data)

			response, err := app.ForgotPassword(
				map[string]string{
					"email": data.Email,
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode())
		},
	)

	t.Run(
		"Returns 400 when request is invalid",
		func(t *testing.T) {
			t.Parallel()
			app, data := setup(t)
			app.CreateAndVerify(t, data)

			testCases := []struct {
				name               string
				data               map[string]string
				expectedErrorField string
			}{
				{
					name: "empty email",
					data: map[string]string{
						"email": "",
					},
					expectedErrorField: "email",
				},
				{
					name: "invalid email",
					data: map[string]string{
						"email": "invalid email",
					},
					expectedErrorField: "email",
				},
			}

			for _, tc := range testCases {
				tc := tc

				t.Run(
					tc.name,
					func(t *testing.T) {
						t.Parallel()

						response, err := app.ForgotPassword(tc.data)

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
			app.CreateAndVerify(t, data)

			err := app.BanUser(data.Email)
			assert.NoError(t, err)

			response, err := app.ForgotPassword(
				map[string]string{
					"email": data.Email,
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		},
	)

	t.Run(
		"Returns 400 when user is not verified",
		func(t *testing.T) {
			t.Parallel()
			app, data := setup(t)

			response, err := app.SignUp(data)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, response.StatusCode())

			forgotPasswordResponse, err := app.ForgotPassword(
				map[string]string{
					"email": data.Email,
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, forgotPasswordResponse.StatusCode())
		},
	)
}
