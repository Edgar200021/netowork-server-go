package repository

import (
	"context"

	"github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"
	. "github.com/Edgar200021/netowork-server-go/.gen/netowork/public/table"
)

func (r *UserRepository) Create(ctx context.Context, user model.Users) (string, error) {
	var model model.Users

	insertStmt := Users.INSERT(
		Users.Email, Users.FirstName, Users.LastName, Users.Password,
		Users.Role,
	).VALUES(
		user.Email, user.FirstName, user.LastName, user.Password,
		user.Role,
	).RETURNING(Users.ID)

	if err := insertStmt.QueryContext(ctx, r.db, &model); err != nil {
		return "", err
	}

	return model.ID.String(), nil
}
