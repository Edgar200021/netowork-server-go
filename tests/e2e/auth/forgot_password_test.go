package auth

import (
	"net/http"
	"testing"

	"github.com/Edgar200021/netowork-server-go/tests/testapp"
	"github.com/alecthomas/assert/v2"
)

func setupForgotPasswordTest(t *testing.T) (*testapp.TestApp, testapp.SignupData) {
	t.Helper()

	return testapp.New(t), testapp.GenerateFakeData[testapp.SignupData](t)

}

func TestForgotPassword_Returns200_When_RequestIsValid(t *testing.T) {
	t.Parallel()
	app, data := setupForgotPasswordTest(t)
	app.CreateAndVerify(t, data)

	response, err := app.ForgotPassword(
		map[string]string{
			"email": data.Email,
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode())
}

func TestForgotPassword_Returns400_When_RequestIsInvalid(t *testing.T) {
	t.Parallel()
	app, data := setupForgotPasswordTest(t)
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

	for _, test := range testCases {
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				response, err := app.ForgotPassword(
					test.data,
				)

				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, response.StatusCode())
				app.AssertValidationErrors(t, response, test.expectedErrorField)
			},
		)
	}
}

func TestForgotPassword_Returns400_When_UserIsBanned(t *testing.T) {
	t.Parallel()
	app, data := setupForgotPasswordTest(t)
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
}

func TestForgotPassword_Returns400_When_UserIsNotVerified(t *testing.T) {
	t.Parallel()
	app, data := setupForgotPasswordTest(t)

	response, err := app.SignUp(data)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode())

	signInResponse, err := app.ForgotPassword(
		map[string]string{
			"email": data.Email,
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, signInResponse.StatusCode())
}
