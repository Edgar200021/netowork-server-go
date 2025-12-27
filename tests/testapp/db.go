package testapp

import (
	"errors"

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
