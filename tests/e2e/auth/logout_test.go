package auth

import (
	"net/http"
	"testing"

	"github.com/Edgar200021/netowork-server-go/pkg/slice_helpers"
	"github.com/Edgar200021/netowork-server-go/tests/testapp"
	"github.com/alecthomas/assert/v2"
)

func setupLogoutTest(t *testing.T) (*testapp.TestApp, testapp.SignupData) {
	t.Helper()

	return testapp.New(t), testapp.GenerateFakeData[testapp.SignupData](t)
}

func TestLogout_Returns200_When_RequestIsValid(t *testing.T) {
	t.Parallel()
	app, data := setupLogoutTest(t)
	cookies := app.CreateAndSignIn(t, data)

	response, err := app.Logout(cookies)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode())

	logoutCookies := slice_helpers.Filter(
		response.Cookies(), func(val *http.Cookie) bool {
			return val.Name == app.Config.SessionCookieName
		},
	)

	assert.Equal(t, 1, len(logoutCookies))
	assert.Equal(t, logoutCookies[0].MaxAge, -1)
}

func TestLogout_Returns400_When_UserIsBanned(t *testing.T) {
	t.Parallel()
	app, data := setupLogoutTest(t)
	cookies := app.CreateAndSignIn(t, data)

	err := app.BanUser(data.Email)
	assert.NoError(t, err)

	response, err := app.Logout(cookies)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode())

}

func TestLogout_Returns401_When_UserIsNotAuthorized(t *testing.T) {
	t.Parallel()
	app, _ := setupLogoutTest(t)

	response, err := app.Logout(nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode())
}
