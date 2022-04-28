package impl

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"main/logging"
	"main/model/entity"
	"main/repository"
	"main/utils"
)

func NewPgOrmAccountRepository(ctx context.Context,
	db *pg.DB) repository.AccountRepository {
	return pgOrmAccountRepository{
		pgOrm: db,
	}
}

type pgOrmAccountRepository struct {
	pgOrm *pg.DB
}

func (p pgOrmAccountRepository) FindByLogin(ctx context.Context, login string) (*entity.Account, error) {
	res := &entity.Account{}
	err := utils.RunWithProfiler(repository.TagGetAccByLogin, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get Account By Login transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		err = tx.Model(res).Where("login = ?0", login).Select()
		if err != nil {
			logging.ErrorFormat("Error selecting account by login: %s", err)
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

func (p pgOrmAccountRepository) FindAll(ctx context.Context) ([]*entity.Account, error) {
	var res []*entity.Account
	err := utils.RunWithProfiler(repository.TagGetAllAccs, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get All Accounts transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		err = tx.Model(&res).Select()
		if err != nil {
			logging.ErrorFormat("Error selecting all accounts: %s", err)
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

func (p pgOrmAccountRepository) FindById(ctx context.Context, accountId int64) (*entity.Account, error) {
	res := &entity.Account{}
	err := utils.RunWithProfiler(repository.TagGetAccById, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get Account By Id transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		err = tx.Model(res).Where("id = ?0", accountId).Select()
		if err != nil {
			logging.ErrorFormat("Error selecting accounts by id: %s", err)
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

func (p pgOrmAccountRepository) Insert(ctx context.Context, account entity.AccountData) (int64, int64, error) {
	acc := entity.Account{
		Id:       0,
		Login:    account.Login,
		Password: account.Password,
		Role:     account.Role,
	}
	if !acc.InvalidRole(string(acc.Role)) {
		err := utils.RunWithProfiler(repository.TagInsAcc, func() error {
			tx, err := p.pgOrm.Begin()
			if err != nil {
				logging.ErrorFormat("Cannot open Insert transaction: %s", err)
				return err
			}
			defer tx.Rollback()
			newProfile := entity.Profile{}
			_, err = tx.Model(&newProfile).Returning("id").Insert()
			if err != nil {
				logging.ErrorFormat("Cannot Insert profile connected with account %v+: %s", account,
					err.Error())
				return err
			}
			acc.ProfileID = newProfile.OrganizationId
			_, err = tx.Model(&acc).Returning("id").Insert()
			if err != nil {
				logging.ErrorFormat("Cannot Insert account %v+: %s", account,
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
			return -1, -1, err
		}
		return acc.ProfileID, acc.Id, nil
	} else {
		logging.ErrorFormat("Invalid Role %s", acc.Role)
		return -1, -1, fmt.Errorf("invalid role %s", acc.Role)
	}

}

func (p pgOrmAccountRepository) Delete(ctx context.Context, accountId int64) error {
	err := utils.RunWithProfiler(repository.TagDelAcc, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Delete transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		acc := &entity.Account{}
		err = tx.Model(&acc).Where("id = ?0", accountId).Select()
		if err != nil {
			logging.ErrorFormat("Error deleting account %d: %s", accountId, err)
			return err
		}

		_, err = tx.Model(&entity.Profile{}).Where("id = ?0", acc.ProfileID).Delete()
		if err != nil {
			logging.ErrorFormat("Error deleting account %d: %s", accountId, err)
			return err
		}

		_, err = tx.Model(&entity.Account{}).Where("id = ?0", accountId).Delete()
		if err != nil {
			logging.ErrorFormat("Error deleting account %d: %s", accountId, err)
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
