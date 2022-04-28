package impl

import (
	"context"
	"github.com/go-pg/pg/v10"
	"main/logging"
	"main/model/entity"
	"main/repository"
	"main/utils"
)

func NewPgOrmProfileRepository(ctx context.Context,
	db *pg.DB) repository.ProfileRepository {
	return pgOrmProfileRepository{
		pgOrm: db,
	}
}

type pgOrmProfileRepository struct {
	pgOrm *pg.DB
}

func (p pgOrmProfileRepository) InsertProfile(ctx context.Context, data *entity.ProfileData) (int64, error) {
	profile := entity.Profile{}
	err := utils.RunWithProfiler(repository.TagInsPrf, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Insert transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		profile = entity.Profile{
			OrganizationId:   0,
			OrganizationName: data.OrganizationName,
			OrganizationType: data.OrganizationType,
		}

		_, err = tx.Model(&profile).Returning("id").Insert()
		if err != nil {
			logging.ErrorFormat("Cannot Insert profile %d: %s", profile.OrganizationId,
				err.Error())
			return err
		}

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return -1, err
	}
	return profile.OrganizationId, nil
}

func (p pgOrmProfileRepository) UpdateProfile(ctx context.Context, id int64, data *entity.ProfileData) error {
	err := utils.RunWithProfiler(repository.TagUpdPrf, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Update transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		profile := entity.Profile{
			OrganizationId:   id,
			OrganizationName: data.OrganizationName,
			OrganizationType: data.OrganizationType,
		}
		_, err = tx.Model(&profile).Where("id = ?0", id).Update()
		if err != nil {
			logging.ErrorFormat("Error updating profile %d: %s", id, err)
			return err
		}

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (p pgOrmProfileRepository) DeleteProfile(ctx context.Context, id int64) error {
	err := utils.RunWithProfiler(repository.TagDelPrf, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Delete transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		_, err = tx.Model(&entity.Profile{}).Where("id = ?0", id).Delete()
		if err != nil {
			logging.ErrorFormat("Error deleting profile %d: %s", id, err)
			return err
		}

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (p pgOrmProfileRepository) GetProfileById(ctx context.Context, id int64) (*entity.Profile, error) {
	res := &entity.Profile{}
	err := utils.RunWithProfiler(repository.TagGetPrfById, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get Profile By Id transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		err = tx.Model(res).Where("id = ?0", id).Select()
		if err != nil {
			logging.ErrorFormat("Error selecting profile by id: %s", err)
			return err
		}

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
