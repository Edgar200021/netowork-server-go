package repository

import (
	"context"

	"github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"
	. "github.com/Edgar200021/netowork-server-go/.gen/netowork/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func (r *UserRepository) Update(
	ctx context.Context, id uuid.UUID, user *model.Users,
) error {
	updateStmt := Users.UPDATE(Users.MutableColumns).MODEL(user).WHERE(Users.ID.EQ(UUID(id)))

	if _, err := updateStmt.ExecContext(ctx, r.db); err != nil {
		return err
	}

	return nil
}
