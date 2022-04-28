package repository

import (
	"context"
	"main/model/entity"
)

const (
	TagInsPrf     = "INSERT PROFILE"
	TagUpdPrf     = "UPDATE PROFILE"
	TagDelPrf     = "DELETE PROFILE"
	TagGetPrfById = "GET PROFILE BY ID"
)

type ProfileRepository interface {
	InsertProfile(ctx context.Context, data *entity.ProfileData) (int64, error)
	UpdateProfile(ctx context.Context, id int64, data *entity.ProfileData) error
	DeleteProfile(ctx context.Context, id int64) error
	GetProfileById(ctx context.Context, id int64) (*entity.Profile, error)
}
