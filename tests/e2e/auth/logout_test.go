package auth

import (
	"net/http"
	"testing"

	"github.com/Edgar200021/netowork-server-go/pkg/slice_helpers"
	"github.com/Edgar200021/netowork-server-go/tests/testapp"
	"github.com/alecthomas/assert/v2"
)

func TestLogout(t *testing.T) {
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
			cookies := app.CreateAndSignIn(t, data)

			response, err := app.Logout(cookies)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode())

			logoutCookies := slice_helpers.Filter(
				response.Cookies(),
				func(val *http.Cookie) bool {
					return val.Name == app.Config.SessionCookieName
				},
			)

			assert.Equal(t, 1, len(logoutCookies))
			assert.Equal(t, -1, logoutCookies[0].MaxAge)
		},
	)

	t.Run(
		"Returns 400 when user is banned",
		func(t *testing.T) {
			t.Parallel()
			app, data := setup(t)
			cookies := app.CreateAndSignIn(t, data)

			err := app.BanUser(data.Email)
			assert.NoError(t, err)

			response, err := app.Logout(cookies)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, response.StatusCode())
		},
	)

	t.Run(
		"Returns 401 when user is not authorized",
		func(t *testing.T) {
			t.Parallel()
			app, _ := setup(t)

			response, err := app.Logout(nil)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, response.StatusCode())
		},
	)
}
