package mapper

import (
	"context"

	"github.com/ONSdigital/dp-frontend-articles-controller/config"
)

type Bulletin struct {
	// TODO define
	Name string `json:"name"`
}

type BulletinModel struct {
	// TODO define
	Name string `json:"model-name"`
}

func Blank(ctx context.Context, bulletin Bulletin, cfg config.Config) BulletinModel {
	var model BulletinModel
	model.Name = bulletin.Name
	return model
}
