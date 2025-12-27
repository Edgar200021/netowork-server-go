package auth

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"
	. "github.com/Edgar200021/netowork-server-go/.gen/netowork/public/table"
	"github.com/Edgar200021/netowork-server-go/tests/testapp"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/stretchr/testify/suite"
)

type signupTestSuite struct {
	suite.Suite
	testApp    *testapp.TestApp
	signupData testapp.SignupData
}

func (s *signupTestSuite) SetupTest() {
	s.signupData = testapp.GenerateFakeData[testapp.SignupData]()
	s.testApp = testapp.New(s.Suite.T())
}

func (s *signupTestSuite) TestSignup_Returns201_When_RequestIsValid() {
	response, err := s.testApp.SignUp(s.signupData)

	s.Require().NoError(err)
	s.Equal(http.StatusCreated, response.StatusCode())
}

func (s *signupTestSuite) TestSignup_ShouldBeSavedIntoDatabase_When_RequestIsValid() {
	response, err := s.testApp.SignUp(s.signupData)

	s.Require().NoError(err)
	s.Equal(http.StatusCreated, response.StatusCode())

	var dbUser model.Users

	stmt := SELECT(
		Users.ID, Users.IsVerified, Users.Password,
	).FROM(Users).WHERE(Users.Email.EQ(Text(s.signupData.Email)))

	err = stmt.Query(s.testApp.Db, &dbUser)

	s.Require().NoError(err)
	s.Require().NotEmpty(dbUser.ID)
	s.Require().NotEqual(s.signupData.Password, dbUser.Password)
	s.Require().False(dbUser.IsVerified)
}

func (s *signupTestSuite) TestSignup_Returns400_When_RequestIsInvalidValid() {
	testCases := []struct {
		name               string
		data               testapp.SignupData
		expectedErrorField string
	}{
		{
			name: "empty email",
			data: func() testapp.SignupData {
				d := s.signupData
				d.Email = ""
				return d
			}(),
			expectedErrorField: "email",
		},
		{
			name: "invalid email",
			data: func() testapp.SignupData {
				d := s.signupData
				d.Email = "invalid"
				return d
			}(),
			expectedErrorField: "email",
		},
		{
			name: "empty password",
			data: func() testapp.SignupData {
				d := s.signupData
				d.Password = ""
				return d
			}(),
			expectedErrorField: "password",
		},
		{
			name: "password too short",
			data: func() testapp.SignupData {
				d := s.signupData
				d.Password = "pass"
				return d
			}(),
			expectedErrorField: "password",
		},
		{
			name: "password too long",
			data: func() testapp.SignupData {
				d := s.signupData
				d.Password = strings.Repeat("q", 41)
				return d
			}(),
			expectedErrorField: "password",
		},
		{
			name: "empty firstName",
			data: func() testapp.SignupData {
				d := s.signupData
				d.FirstName = ""
				return d
			}(),
			expectedErrorField: "firstName",
		},
		{
			name: "firstName too short",
			data: func() testapp.SignupData {
				d := s.signupData
				d.FirstName = "qw"
				return d
			}(),
			expectedErrorField: "firstName",
		},
		{
			name: "firstName too long",
			data: func() testapp.SignupData {
				d := s.signupData
				d.FirstName = strings.Repeat("q", 41)
				return d
			}(),
			expectedErrorField: "firstName",
		},
		{
			name: "empty lastName",
			data: func() testapp.SignupData {
				d := s.signupData
				d.LastName = ""
				return d
			}(),
			expectedErrorField: "lastName",
		},
		{
			name: "lastName too short",
			data: func() testapp.SignupData {
				d := s.signupData
				d.LastName = "qw"
				return d
			}(),
			expectedErrorField: "lastName",
		},
		{
			name: "lastName too long",
			data: func() testapp.SignupData {
				d := s.signupData
				d.LastName = strings.Repeat("q", 41)
				return d
			}(),
			expectedErrorField: "lastName",
		},
		{
			name: "empty role",
			data: func() testapp.SignupData {
				d := s.signupData
				d.Role = ""
				return d
			}(),
			expectedErrorField: "role",
		},
		{
			name: "invalid role",
			data: func() testapp.SignupData {
				d := s.signupData
				d.Role = "invalidrole"
				return d
			}(),
			expectedErrorField: "role",
		},
		{
			name: "admin role",
			data: func() testapp.SignupData {
				d := s.signupData
				d.Role = "admin"
				return d
			}(),
			expectedErrorField: "role",
		},
	}

	for _, test := range testCases {
		s.Run(
			test.name, func() {
				response, err := s.testApp.SignUp(test.data)
				s.Require().NoError(err)

				s.Require().Equal(http.StatusBadRequest, response.StatusCode())
				s.testApp.AssertValidationErrors(s.T(), response, test.expectedErrorField)
			},
		)
	}
}

func (s *signupTestSuite) TestSignup_Return400_When_UserAlreadyExists() {
	response, err := s.testApp.SignUp(s.signupData)

	s.Require().NoError(err)
	s.Equal(http.StatusCreated, response.StatusCode())

	secondResponse, err := s.testApp.SignUp(s.signupData)

	s.Require().NoError(err)
	s.Equal(http.StatusBadRequest, secondResponse.StatusCode())
}

func TestSignupTestSuite(t *testing.T) {
	suite.Run(t, new(signupTestSuite))
}
