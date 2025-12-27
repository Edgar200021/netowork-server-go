package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Edgar200021/netowork-server-go/tests/testapp"
	"github.com/alecthomas/assert/v2"
	"github.com/go-faker/faker/v4"
)

func setupResetPasswordTest(t *testing.T) (*testapp.TestApp, testapp.SignupData) {
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

func TestResetPassword_Returns200_When_RequestIsValid(t *testing.T) {
	t.Parallel()
	app, data := setupResetPasswordTest(t)

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
}

func TestResetPassword_ShouldApplyChangesIntoDatabase_When_RequestIsValid(t *testing.T) {
	t.Parallel()
	app, data := setupResetPasswordTest(t)

	token, err := app.GetResetPasswordToken()

	assert.NoError(t, err)
	assert.NotEqual(t, "", token)

	user, err := app.GetUser(data.Email)
	assert.NoError(t, err)
	assert.NotZero(t, user)

	response, err := app.ResetPassword(
		map[string]string{
			"email":    data.Email,
			"token":    token,
			"password": strings.Repeat("q", testapp.PasswordMinLength),
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode())

	userAfterChanges, err := app.GetUser(data.Email)

	assert.NoError(t, err)
	assert.True(t, userAfterChanges.IsVerified)
	assert.NotEqual(t, user.Password, userAfterChanges.Password)
}

func TestResetPassword_ShouldApplyChangesIntoDatabase_When_RequestIsInvalid(t *testing.T) {
	t.Parallel()
	app, _ := setupResetPasswordTest(t)

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
				"email":    "",
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

	for _, test := range testCases {
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				response, err := app.ResetPassword(test.data)

				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, response.StatusCode())
			},
		)
	}
}

func TestResetPassword_Returns400_When_UserIsBanned(t *testing.T) {
	t.Parallel()
	app, data := setupResetPasswordTest(t)

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
}

func TestResetPassword_Returns400_When_EmailIsDifferent(t *testing.T) {
	t.Parallel()
	app, _ := setupResetPasswordTest(t)

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
}
