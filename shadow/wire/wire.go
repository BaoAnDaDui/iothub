//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"gorm.io/gorm"
	"iothub/shadow"

	"github.com/google/wire"
)

func InitSvc(dbConn *gorm.DB, conn shadow.StatusGetter) shadow.Service {
	wire.Build(
		shadow.NewSvc,
		shadow.NewShadowRepo,
	)
	return nil
}
