package http

import (
	"github.com/google/wire"
	"member-pre/internal/interfaces/http/handler"
)

var WireHttpSet = wire.NewSet(
	handler.NewAuthHandler,
	handler.NewUserHandler,
	handler.NewStoreHandler,
	handler.NewSlotTemplateHandler,
	handler.NewSlotHandler,
	NewAppRouteRegistrar,
	// 注意: NewUserHandler现在需要slotService参数
)
