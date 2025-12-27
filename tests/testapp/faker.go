package testapp

import "github.com/go-faker/faker/v4"

type SignupData struct {
	FirstName string `faker:"first_name,len=10"`
	LastName  string `faker:"first_name,len=10"`
	Email     string `faker:"email,unique"`
	Password  string `faker:"password, len=20"`
	Role      string `faker:"oneof: client, freelancer"`
}

func GenerateFakeData[T any]() T {
	var data T

	err := faker.FakeData(&data)
	if err != nil {
		panic(err)
	}

	return data
}
