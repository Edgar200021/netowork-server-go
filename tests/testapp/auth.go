package testapp

import (
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
