package repo

import (
	"context"

	"github.com/suzmii/ACMBot/internal/database/dbmodel"
	"github.com/suzmii/ACMBot/internal/database/gen"
)

func SaveUser(ctx context.Context, user *dbmodel.CodeforcesUser) error {
	err := gen.Q.WithContext(ctx).
		CodeforcesUser.
		Save(user)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(ctx context.Context, username string) (*dbmodel.CodeforcesUser, error) {
	user, err := gen.Q.WithContext(ctx).
		CodeforcesUser.
		Where(gen.CodeforcesUser.Username.Eq(username)).
		Take()
	if err != nil {
		return nil, err
	}
	return user, nil
}
