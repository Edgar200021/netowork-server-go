package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Edgar200021/netowork-server-go/pkg/slice_helpers"
	"github.com/Edgar200021/netowork-server-go/tests/testapp"
	"github.com/alecthomas/assert/v2"
	"github.com/go-faker/faker/v4"
)

func setupSignInTest(t *testing.T) (*testapp.TestApp, testapp.SignupData) {
	t.Helper()

	return testapp.New(t), testapp.GenerateFakeData[testapp.SignupData](t)
}

func TestSignIn_Returns200_When_RequestIsValid(t *testing.T) {
	t.Parallel()
	app, data := setupSignInTest(t)
	app.CreateAndVerify(t, data)

	response, err := app.SignIn(
		map[string]string{
			"email":    data.Email,
			"password": data.Password,
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode())

	cookies := slice_helpers.Filter(
		response.Cookies(), func(val *http.Cookie) bool {
			return val.Name == app.Config.SessionCookieName
		},
	)

	assert.Equal(t, 1, len(cookies))
	assert.True(t, true, cookies[0].HttpOnly)
	assert.True(t, true, cookies[0].HttpOnly)
}

func TestSignIn_Returns400_When_RequestIsInvalid(t *testing.T) {
	t.Parallel()
	app, data := setupSignInTest(t)

	testCases := []struct {
		name               string
		data               map[string]string
		expectedErrorField string
	}{
		{
			name: "empty email",
			data: map[string]string{
				"password": data.Password,
			},
			expectedErrorField: "email",
		},
		{
			name: "invalid email",
			data: map[string]string{
				"email":    "invalid email",
				"password": data.Password,
			},
			expectedErrorField: "email",
		},
		{
			name: "empty password",
			data: map[string]string{
				"email": data.Email,
			},
			expectedErrorField: "password",
		},
		{
			name: "password too short",
			data: map[string]string{
				"email":    data.Email,
				"password": strings.Repeat("q", testapp.PasswordMinLength-1),
			},
			expectedErrorField: "password",
		},
		{
			name: "password too long",
			data: map[string]string{
				"email":    data.Email,
				"password": strings.Repeat("q", testapp.PasswordMaxLength+1),
			},
			expectedErrorField: "password",
		},
	}

	for _, test := range testCases {
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				response, err := app.SignIn(test.data)

				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, response.StatusCode())
				app.AssertValidationErrors(t, response, test.expectedErrorField)
			},
		)
	}
}

func TestSignIn_Returns400_When_UserIsBanned(t *testing.T) {
	t.Parallel()
	app, data := setupSignInTest(t)
	app.CreateAndVerify(t, data)

	err := app.BanUser(data.Email)
	assert.NoError(t, err)

	response, err := app.SignIn(
		map[string]string{
			"email":    data.Email,
			"password": data.Password,
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode())
}

func TestSignIn_Returns400_When_UserIsNotVerified(t *testing.T) {
	t.Parallel()
	app, data := setupSignInTest(t)

	response, err := app.SignUp(data)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode())

	signInResponse, err := app.SignIn(
		map[string]string{
			"email":    data.Email,
			"password": data.Password,
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, signInResponse.StatusCode())
}

func TestSignIn_Return400_When_CredentialsAreInvalid(t *testing.T) {
	t.Parallel()
	app, data := setupSignInTest(t)

	testCases := []struct {
		name string
		data map[string]string
	}{
		{
			name: "incorrect email",
			data: map[string]string{
				"email":    faker.Email(),
				"password": data.Password,
			},
		},
		{
			name: "incorrect password",
			data: map[string]string{
				"email":    data.Email,
				"password": strings.Repeat("q", testapp.PasswordMinLength+1),
			},
		},
	}

	for _, test := range testCases {
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()

				response, err := app.SignIn(test.data)

				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, response.StatusCode())
			},
		)
	}

}
