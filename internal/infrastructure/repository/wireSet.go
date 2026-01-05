package repository

import (
	"github.com/google/wire"
	"member-pre/internal/domain/store"
	"member-pre/internal/domain/user"
)

var WireRepoSet = wire.NewSet(
	NewUserRepository,
	NewStoreRepository,
	// 绑定接口和实现
	wire.Bind(new(user.IUserRepository), new(*UserRepository)),
	wire.Bind(new(store.IStoreRepository), new(*StoreRepository)),
)
