package impl

import (
	"context"
	"github.com/go-pg/pg/v10"
	"main/logging"
	"main/model/entity"
	"main/repository"
	"main/utils"
)

func NewPgOrmProductRepository(ctx context.Context,
	db *pg.DB) repository.ProductRepository {
	return pgOrmProductRepository{
		pgOrm: db,
	}
}

type pgOrmProductRepository struct {
	pgOrm *pg.DB
}

func (p pgOrmProductRepository) InsertProducts(ctx context.Context, products []*entity.Product) error {
	err := utils.RunWithProfiler(repository.TagInsPr, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Insert transaction: %s", err)
			return err
		}
		defer tx.Rollback()
		for _, product := range products {
			//TODO handle returned id
			_, err = tx.Model(product).Returning("id").Insert()
			if err != nil {
				logging.ErrorFormat("Cannot Insert product %d: %s", product.Id,
					err.Error())
				return err
			}
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

func (p pgOrmProductRepository) UpdateProduct(ctx context.Context, id int, data entity.ProductData) error {
	err := utils.RunWithProfiler(repository.TagUpdPr, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Update transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		product := entity.Product{
			Id:          int64(id),
			Description: data.Description,
			CategoryId:  data.CategoryId,
			Name:        data.Name,
		}
		_, err = tx.Model(&product).Where("id = ?0", id).Update()
		if err != nil {
			logging.ErrorFormat("Error updating product %d: %s", id, err)
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

func (p pgOrmProductRepository) DeleteProducts(ctx context.Context, id int) error {
	err := utils.RunWithProfiler(repository.TagDelPr, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Delete transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		_, err = tx.Model(&entity.Product{}).Where("id = ?0", id).Delete()
		if err != nil {
			logging.ErrorFormat("Error deleting product %d: %s", id, err)
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

func (p pgOrmProductRepository) GetAllProducts(ctx context.Context) ([]*entity.Product, error) {
	var res []*entity.Product
	err := utils.RunWithProfiler(repository.TagGetAllPr, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get All Products transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		err = tx.Model(&res).Relation("Category").Select()
		if err != nil {
			logging.ErrorFormat("Error selecting all products: %s", err)
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

func (p pgOrmProductRepository) GetProductById(ctx context.Context, id int) (*entity.Product, error) {
	res := &entity.Product{}
	err := utils.RunWithProfiler(repository.TagGetPrById, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get Product By Id transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		err = tx.Model(res).Relation("Category").Where("product.id = ?0", id).Select()
		if err != nil {
			logging.ErrorFormat("Error selecting product by id: %s", err)
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

func (p pgOrmProductRepository) GetProductsByQuery(ctx context.Context, q string) ([]*entity.Product, error) {
	var res []*entity.Product
	err := utils.RunWithProfiler(repository.TagGetPrByQ, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get Products By Query transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		err = tx.Model(&res).Relation("Category").
			Where("product.description LIKE ?0 OR category.category_name LIKE ?1",
				"%"+q+"%", "%"+q+"%").Select()
		if err != nil {
			logging.ErrorFormat("Error selecting products: %s", err)
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
