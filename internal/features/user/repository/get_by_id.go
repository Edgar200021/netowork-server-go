package repository

import (
	"context"
	"errors"

	"github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"
	. "github.com/Edgar200021/netowork-server-go/.gen/netowork/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

func (r *UserRepository) GetById(ctx context.Context, id uuid.UUID) (*model.Users, error) {
	var user model.Users

	query := SELECT(Users.AllColumns).FROM(Users).WHERE(Users.ID.EQ(UUID(id)))

	if err := query.QueryContext(ctx, r.db, &user); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
