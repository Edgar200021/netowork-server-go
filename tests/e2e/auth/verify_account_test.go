package auth

import (
	"net/http"
	"testing"

	"github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"
	. "github.com/Edgar200021/netowork-server-go/.gen/netowork/public/table"
	"github.com/Edgar200021/netowork-server-go/tests/testapp"
	"github.com/alecthomas/assert/v2"
	"github.com/go-faker/faker/v4"
	. "github.com/go-jet/jet/v2/postgres"
)

func setupVerifyAccountTest(t *testing.T) (*testapp.TestApp, testapp.SignupData, string) {
	t.Helper()

	signupData := testapp.GenerateFakeData[testapp.SignupData](t)
	app := testapp.New(t)

	response, err := app.SignUp(signupData)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode())

	verificationToken, err := app.GetVerificationToken()

	assert.NoError(t, err)
	assert.NotEqual(t, verificationToken, "")

	return app, signupData, verificationToken
}

func TestVerifyAccount_Returns200_When_RequestIsValid(t *testing.T) {
	t.Parallel()
	app, data, token := setupVerifyAccountTest(t)

	verifyAccountResponse, err := app.VerifyAccount(
		map[string]string{
			"email": data.Email,
			"token": token,
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, verifyAccountResponse.StatusCode())
}

func TestVerifyAccount_ShouldApplyChangesIntoDatabase_When_RequestIsValid(t *testing.T) {
	t.Parallel()
	app, data, token := setupVerifyAccountTest(t)

	verifyAccountResponse, err := app.VerifyAccount(
		map[string]string{
			"email": data.Email,
			"token": token,
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, verifyAccountResponse.StatusCode())

	var dbUser model.Users

	stmt := SELECT(Users.IsVerified).FROM(Users).WHERE(Users.Email.EQ(Text(data.Email)))

	err = stmt.Query(app.Db, &dbUser)

	assert.NoError(t, err)
	assert.True(t, dbUser.IsVerified)
}

func TestVerifyAccount_Return400_When_RequestIsInvalid(t *testing.T) {
	t.Parallel()
	app, _, _ := setupVerifyAccountTest(t)

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

	for _, test := range testCases {
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()
				response, err := app.VerifyAccount(test.data)

				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, response.StatusCode())
				app.AssertValidationErrors(t, response, test.expectedErrorField)
			},
		)
	}
}

func TestVerifyAccount_Returns400_When_TokenIsInvalidOrEmailIsDifferent(t *testing.T) {
	t.Parallel()
	app, data, token := setupVerifyAccountTest(t)

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
			name: "random token",
			data: map[string]string{
				"email": faker.Email(),
				"token": token,
			},
		},
	}

	for _, test := range testCases {
		t.Run(
			test.name, func(t *testing.T) {

				t.Parallel()
				response, err := app.VerifyAccount(test.data)

				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, response.StatusCode())
			},
		)
	}

}

func TestVerifyAccount_Returns400_When_UserIsBanned(t *testing.T) {
	t.Parallel()
	app, data, token := setupVerifyAccountTest(t)

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
}

func TestVerifyAccount_Returns400_When_TokenExpired(t *testing.T) {
	t.Parallel()
	app, data, token := setupVerifyAccountTest(t)

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
}

func TestVerifyAccount_Returns404_When_UserIsNotFound(t *testing.T) {
	t.Parallel()
	app, data, token := setupVerifyAccountTest(t)

	deleteStmt := Users.DELETE().WHERE(Users.Email.EQ(Text(data.Email)))
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
}
