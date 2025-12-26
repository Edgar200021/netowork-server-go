package repository

import (
	"context"
	"errors"

	"github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"
	. "github.com/Edgar200021/netowork-server-go/.gen/netowork/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.Users, error) {
	var user model.Users

	stmt := SELECT(Users.AllColumns).FROM(Users).WHERE(Users.Email.EQ(Text(email)))

	if err := stmt.QueryContext(ctx, r.db, &user); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
