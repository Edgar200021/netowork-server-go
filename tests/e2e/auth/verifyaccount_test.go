package auth

import (
	"net/http"
	"testing"

	"github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"
	. "github.com/Edgar200021/netowork-server-go/.gen/netowork/public/table"
	"github.com/Edgar200021/netowork-server-go/tests/testapp"
	"github.com/go-faker/faker/v4"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/stretchr/testify/suite"
)

type verifyAccountTestSuite struct {
	suite.Suite
	testApp           *testapp.TestApp
	signupData        testapp.SignupData
	verificationToken string
}

func (s *verifyAccountTestSuite) SetupTest() {
	s.signupData = testapp.GenerateFakeData[testapp.SignupData]()
	s.testApp = testapp.New(s.Suite.T())

	response, err := s.testApp.SignUp(s.signupData)

	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, response.StatusCode())

	verificationToken, err := s.testApp.GetVerificationToken()

	s.Require().NoError(err)
	s.Require().NotEmpty(verificationToken)

	s.verificationToken = verificationToken
}

func (s *verifyAccountTestSuite) TestVerifyAccount_Returns200_When_RequestIsValid() {
	verifyAccountResponse, err := s.testApp.VerifyAccount(
		map[string]string{
			"email": s.signupData.Email,
			"token": s.verificationToken,
		},
	)

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, verifyAccountResponse.StatusCode())
}

func (s *verifyAccountTestSuite) TestVerifyAccount_ShouldApplyChangesIntoDatabase_When_RequestIsValid() {
	verifyAccountResponse, err := s.testApp.VerifyAccount(
		map[string]string{
			"email": s.signupData.Email,
			"token": s.verificationToken,
		},
	)

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, verifyAccountResponse.StatusCode())

	var dbUser model.Users

	stmt := SELECT(Users.IsVerified).FROM(Users).WHERE(Users.Email.EQ(Text(s.signupData.Email)))

	err = stmt.Query(s.testApp.Db, &dbUser)

	s.Require().NoError(err)
	s.Require().True(dbUser.IsVerified)
}

func (s *verifyAccountTestSuite) TestVerifyAccount_Return400_When_RequestIsInvalid() {
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
		s.Run(
			test.name, func() {
				response, err := s.testApp.VerifyAccount(test.data)

				s.Require().NoError(err)
				s.Require().Equal(http.StatusBadRequest, response.StatusCode())
				s.testApp.AssertValidationErrors(s.T(), response, test.expectedErrorField)
			},
		)
	}
}

func (s *verifyAccountTestSuite) TestVerifyAccount_Returns400_When_TokenIsInvalidOrEmailIsDifferent() {
	testCases := []struct {
		name string
		data map[string]string
	}{
		{
			name: "random token",
			data: map[string]string{
				"email": s.signupData.Email,
				"token": "random token",
			},
		},
		{
			name: "random token",
			data: map[string]string{
				"email": faker.Email(),
				"token": s.verificationToken,
			},
		},
	}

	for _, test := range testCases {
		s.Run(
			test.name, func() {
				response, err := s.testApp.VerifyAccount(test.data)

				s.Require().NoError(err)
				s.Require().Equal(http.StatusBadRequest, response.StatusCode())
			},
		)
	}

}

func (s *verifyAccountTestSuite) TestVerifyAccount_Returns400_When_UserIsBanned() {
	updateStmt := Users.UPDATE(Users.IsBanned).SET(Users.IsBanned.SET(Bool(true))).WHERE(
		Users.Email.EQ(
			Text(
				s.signupData.
					Email,
			),
		),
	)
	res, err := updateStmt.Exec(s.testApp.Db)
	s.Require().NoError(err)

	updatedRows, err := res.RowsAffected()

	s.Require().NoError(err)
	s.Require().EqualValues(1, updatedRows)

	response, err := s.testApp.VerifyAccount(
		map[string]string{
			"email": s.signupData.Email,
			"token": s.verificationToken,
		},
	)

	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, response.StatusCode())
}

func (s *verifyAccountTestSuite) TestVerifyAccount_Returns404_When_UserIsNotFound() {
	deleteStmt := Users.DELETE().WHERE(Users.Email.EQ(Text(s.signupData.Email)))
	res, err := deleteStmt.Exec(s.testApp.Db)

	s.Require().NoError(err)

	deletedRows, err := res.RowsAffected()

	s.Require().NoError(err)
	s.Require().EqualValues(1, deletedRows)

	response, err := s.testApp.VerifyAccount(
		map[string]string{
			"email": s.signupData.Email,
			"token": s.verificationToken,
		},
	)

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, response.StatusCode())
}

func TestVerifyAccountTestSuite(t *testing.T) {
	suite.Run(t, new(verifyAccountTestSuite))
}
