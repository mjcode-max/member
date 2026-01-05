package domain

import (
	"github.com/google/wire"
	"member-pre/internal/domain/auth"
	"member-pre/internal/domain/store"
	"member-pre/internal/domain/user"
)

var WireDoMainSet = wire.NewSet(
	user.NewUserService,
	auth.NewAuthService,
	store.NewStoreService,
	// 绑定接口
	wire.Bind(new(auth.IUserService), new(*user.UserService)),
)
