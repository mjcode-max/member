package http

import (
	"github.com/google/wire"
	"member-pre/internal/interfaces/http/handler"
)

var WireHttpSet = wire.NewSet(
	handler.NewAuthHandler,
	handler.NewUserHandler,
	handler.NewStoreHandler,
	NewAppRouteRegistrar,
)
