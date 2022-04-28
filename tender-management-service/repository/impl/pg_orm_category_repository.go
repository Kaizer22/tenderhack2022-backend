package impl

import (
	"context"
	"github.com/go-pg/pg/v10"
	"main/logging"
	"main/model/entity"
	"main/repository"
	"main/utils"
)

func NewPgOrmCategoryRepository(ctx context.Context,
	db *pg.DB) repository.CategoryRepository {
	return pgOrmCategoryRepository{
		pgOrm: db,
	}
}

type pgOrmCategoryRepository struct {
	pgOrm *pg.DB
}

func (p pgOrmCategoryRepository) InsertCategories(ctx context.Context, categories []*entity.Category) error {
	err := utils.RunWithProfiler(repository.TagInsCat, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Insert transaction: %s", err)
			return err
		}
		defer tx.Rollback()
		//TODO handle returned id
		_, err = tx.Model(&categories).Returning("id").Insert()
		if err != nil {
			logging.ErrorFormat("Cannot Insert categories %v+: %s", categories,
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
		return err
	}
	return nil
}

func (p pgOrmCategoryRepository) UpdateCategory(ctx context.Context, id int, data entity.CategoryData) error {
	err := utils.RunWithProfiler(repository.TagUpdCat, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Update transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		category := entity.Category{
			Id:          int64(id),
			Description: data.Description,
			Name:        data.Name,
		}
		_, err = tx.Model(&category).Where("id = ?0", id).Update()
		if err != nil {
			logging.ErrorFormat("Error updating category %d: %s", id, err)
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

func (p pgOrmCategoryRepository) DeleteCategory(ctx context.Context, id int) error {
	err := utils.RunWithProfiler(repository.TagDelCat, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Delete transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		_, err = tx.Model(&entity.Category{}).Where("id = ?0", id).Delete()
		if err != nil {
			logging.ErrorFormat("Error deleting category %d: %s", id, err)
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

func (p pgOrmCategoryRepository) GetAllCategories(ctx context.Context) ([]*entity.Category, error) {
	var res []*entity.Category
	err := utils.RunWithProfiler(repository.TagGetAllCats, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get All Categories transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		err = tx.Model(&res).Select()
		if err != nil {
			logging.ErrorFormat("Error selecting all categories: %s", err)
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

func (p pgOrmCategoryRepository) GetCategoryById(ctx context.Context, id int) (*entity.Category, error) {
	res := &entity.Category{}
	err := utils.RunWithProfiler(repository.TagGetCatById, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get Category By Id transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		err = tx.Model(res).Where("id = ?0", id).Select()
		if err != nil {
			logging.ErrorFormat("Error selecting categories by id: %s", err)
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
