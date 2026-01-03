package repository

import (
	"github.com/google/wire"
	"member-pre/internal/domain/user"
)

var WireRepoSet = wire.NewSet(
	NewUserRepository,
	// 绑定接口和实现
	wire.Bind(new(user.IUserRepository), new(*UserRepository)),
)
