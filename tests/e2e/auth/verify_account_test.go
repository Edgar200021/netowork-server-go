package auth

import (
	"net/http"
	"testing"

	. "github.com/Edgar200021/netowork-server-go/.gen/netowork/public/table"
	"github.com/Edgar200021/netowork-server-go/tests/testapp"
	"github.com/alecthomas/assert/v2"
	"github.com/go-faker/faker/v4"
	. "github.com/go-jet/jet/v2/postgres"
)

func TestVerifyAccount(t *testing.T) {
	setup := func(t *testing.T) (*testapp.TestApp, testapp.SignupData, string) {
		t.Helper()

		signupData := testapp.GenerateFakeData[testapp.SignupData](t)
		app := testapp.New(t)

		response, err := app.SignUp(signupData)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.StatusCode())

		verificationToken, err := app.GetVerificationToken()
		assert.NoError(t, err)
		assert.NotEqual(t, "", verificationToken)

		return app, signupData, verificationToken
	}
	t.Parallel()

	t.Run(
		"Returns 200 when request is successful",
		func(t *testing.T) {
			t.Parallel()
			app, data, token := setup(t)

			response, err := app.VerifyAccount(
				map[string]string{
					"email": data.Email,
					"token": token,
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode())
		},
	)

	t.Run(
		"Should apply changes into database when request is successful",
		func(t *testing.T) {
			t.Parallel()
			app, data, token := setup(t)

			response, err := app.VerifyAccount(
				map[string]string{
					"email": data.Email,
					"token": token,
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode())

			dbUser, err := app.GetUser(data.Email)
			assert.NoError(t, err)
			assert.True(t, dbUser.IsVerified)
		},
	)

	t.Run(
		"Returns 400 when request is invalid",
		func(t *testing.T) {
			t.Parallel()
			app, _, _ := setup(t)

			testCases := []struct {
				name               string
				data               map[string]string
				expectedErrorField string
			}{
				{
					name: "empty email",
					data: map[string]string{
						"token": "some token",
					},
					expectedErrorField: "email",
				},
				{
					name: "invalid email",
					data: map[string]string{
						"email": "invalid email",
						"token": "some token",
					},
					expectedErrorField: "email",
				},
				{
					name: "empty token",
					data: map[string]string{
						"email": faker.Email(),
					},
					expectedErrorField: "token",
				},
			}

			for _, tc := range testCases {
				tc := tc

				t.Run(
					tc.name,
					func(t *testing.T) {
						t.Parallel()
						response, err := app.VerifyAccount(tc.data)

						assert.NoError(t, err)
						assert.Equal(t, http.StatusBadRequest, response.StatusCode())
						app.AssertValidationErrors(t, response, tc.expectedErrorField)
					},
				)
			}
		},
	)

	t.Run(
		"Returns 400 when token is invalid or email is different",
		func(t *testing.T) {
			t.Parallel()
			app, data, token := setup(t)

			testCases := []struct {
				name string
				data map[string]string
			}{
				{
					name: "random token",
					data: map[string]string{
						"email": data.Email,
						"token": "random token",
					},
				},
				{
					name: "different email",
					data: map[string]string{
						"email": faker.Email(),
						"token": token,
					},
				},
			}

			for _, tc := range testCases {
				tc := tc

				t.Run(
					tc.name,
					func(t *testing.T) {
						t.Parallel()
						response, err := app.VerifyAccount(tc.data)

						assert.NoError(t, err)
						assert.Equal(t, http.StatusBadRequest, response.StatusCode())
					},
				)
			}
		},
	)

	t.Run(
		"Returns 400 when user is banned",
		func(t *testing.T) {
			t.Parallel()
			app, data, token := setup(t)

			err := app.BanUser(data.Email)
			assert.NoError(t, err)

			response, err := app.VerifyAccount(
				map[string]string{
					"email": data.Email,
					"token": token,
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		},
	)

	t.Run(
		"Returns 400 when token is expired",
		func(t *testing.T) {
			t.Parallel()
			app, data, token := setup(t)

			err := app.ExpireVerificationToken()
			assert.NoError(t, err)

			response, err := app.VerifyAccount(
				map[string]string{
					"email": data.Email,
					"token": token,
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		},
	)

	t.Run(
		"Returns 404 when user is not found",
		func(t *testing.T) {
			t.Parallel()
			app, data, token := setup(t)

			deleteStmt := Users.
				DELETE().
				WHERE(Users.Email.EQ(Text(data.Email)))

			res, err := deleteStmt.Exec(app.Db)
			assert.NoError(t, err)

			deletedRows, err := res.RowsAffected()
			assert.NoError(t, err)
			assert.Equal(t, 1, deletedRows)

			response, err := app.VerifyAccount(
				map[string]string{
					"email": data.Email,
					"token": token,
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, response.StatusCode())
		},
	)
}
