package repository

import (
	"context"
	"main/model/entity"
)

const (
	TagInsAcc        = "INSERT ACCOUNT"
	TagDelAcc        = "DELETE ACCOUNT"
	TagGetAllAccs    = "GET ALL ACCOUNTS"
	TagGetAccByLogin = "GET ACCOUNT BY LOGIN"
	TagGetAccById    = "GET ACCOUNT BY ID"
)

type AccountRepository interface {
	FindAll(ctx context.Context) ([]*entity.Account, error)
	FindById(ctx context.Context, accountId int64) (*entity.Account, error)
	FindByLogin(ctx context.Context, login string) (*entity.Account, error)
	Insert(ctx context.Context, account entity.AccountData) (int64, int64, error)
	Delete(ctx context.Context, accountId int64) error
}
