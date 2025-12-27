package testapp

import (
	"net/http"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/go-resty/resty/v2"
)

func (a *TestApp) SignUp(data SignupData) (*resty.Response, error) {
	response, err := a.client.R().SetBody(data).Post(a.addressV1 + "/auth/sign-up")
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *TestApp) SignIn(data interface{}) (*resty.Response, error) {
	response, err := a.client.R().SetBody(data).Post(a.addressV1 + "/auth/sign-in")
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *TestApp) VerifyAccount(data interface{}) (*resty.Response, error) {
	response, err := a.client.R().SetBody(data).Post(a.addressV1 + "/auth/verify-account")
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *TestApp) ForgotPassword(data interface{}) (
	*resty.Response, error,
) {
	response, err := a.client.R().SetBody(data).Post(
		a.
			addressV1 + "/auth/forgot-password",
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *TestApp) ResetPassword(data interface{}) (
	*resty.Response,
	error,
) {
	response, err := a.client.R().SetBody(data).Post(
		a.
			addressV1 + "/auth/reset-password",
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *TestApp) Logout(cookies []*http.Cookie) (*resty.Response, error) {
	response, err := a.client.R().SetCookies(cookies).Post(a.addressV1 + "/auth/logout")
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *TestApp) CreateAndVerify(t *testing.T, data SignupData) {
	response, err := a.SignUp(data)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode())

	verificationToken, err := a.GetVerificationToken()

	assert.NoError(t, err)
	assert.NotEqual(t, "", verificationToken)

	verificationResponse, err := a.VerifyAccount(
		map[string]string{
			"email": data.Email,
			"token": verificationToken,
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, verificationResponse.StatusCode())
}

func (a *TestApp) CreateAndSignIn(t *testing.T, data SignupData) []*http.Cookie {
	a.CreateAndVerify(t, data)

	signInResponse, err := a.SignIn(
		map[string]string{
			"email":    data.Email,
			"password": data.Password,
		},
	)

	assert.NoError(t, err)
	assert.Equal(
		t, http.StatusOK,
		signInResponse.StatusCode(),
	)

	return signInResponse.Cookies()
}
