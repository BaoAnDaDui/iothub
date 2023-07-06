//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"gorm.io/gorm"
	"iothub/thing"

	"github.com/google/wire"
	"iothub/pkg/uuid"
	"iothub/shadow"
)

func InitSvc(ctx context.Context, dbConn *gorm.DB, shadowSvc shadow.Service, connector shadow.Connectivity) thing.Service {
	wire.Build(
		thing.NewThingRepo,
		uuid.New,
		thing.NewSvc,
	)
	return nil
}
