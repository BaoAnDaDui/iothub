// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"context"
	"gorm.io/gorm"
	"iothub/pkg/uuid"
	"iothub/shadow"
	"iothub/thing"
)

// Injectors from wire.go:

func InitSvc(ctx context.Context, dbConn *gorm.DB, shadowSvc shadow.Service, connector shadow.Connectivity) thing.Service {
	repo := thing.NewThingRepo(dbConn)
	idProvider := uuid.New()
	service := thing.NewSvc(repo, idProvider, shadowSvc, connector)
	return service
}