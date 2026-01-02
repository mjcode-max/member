package repository

import (
	"github.com/google/wire"
	"member-pre/internal/domain/auth"
)

var WireRepoSet = wire.NewSet(
	NewAuthRepository,
	// 绑定接口和实现
	wire.Bind(new(auth.IAuthRepository), new(*AuthRepository)),
)
