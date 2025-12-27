package testapp

import (
	"errors"

	"github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"
	. "github.com/Edgar200021/netowork-server-go/.gen/netowork/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

func (a *TestApp) BanUser(email string) error {
	updateStmt := Users.UPDATE(Users.IsBanned).SET(Users.IsBanned.SET(Bool(true))).WHERE(
		Users.Email.EQ(
			Text(
				email,
			),
		),
	)
	res, err := updateStmt.Exec(a.Db)
	if err != nil {
		return err
	}

	updatedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if updatedRows != 1 {
		return errors.New("failed to ban user")
	}

	return nil

}

func (a *TestApp) GetUser(email string) (*model.Users, error) {
	var user model.Users

	stmt := SELECT(Users.AllColumns).FROM(Users).WHERE(Users.Email.EQ(Text(email)))

	err := stmt.Query(a.Db, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
