package repository

import (
	"github.com/google/wire"
	"member-pre/internal/domain/slot"
	"member-pre/internal/domain/member"
	"member-pre/internal/domain/store"
	"member-pre/internal/domain/user"
)

var WireRepoSet = wire.NewSet(
	NewUserRepository,
	NewStoreRepository,
	NewTemplateRepository,
	NewSlotRepository,
	// 绑定接口和实现
	wire.Bind(new(user.IUserRepository), new(*UserRepository)),
	wire.Bind(new(store.IStoreRepository), new(*StoreRepository)),
	wire.Bind(new(slot.ITemplateRepository), new(*TemplateRepository)),
	wire.Bind(new(slot.ISlotRepository), new(*SlotRepository)),
	NewMemberRepository,
	NewUsageRepository,
	// 绑定接口和实现
	wire.Bind(new(user.IUserRepository), new(*UserRepository)),
	wire.Bind(new(store.IStoreRepository), new(*StoreRepository)),
	wire.Bind(new(member.IMemberRepository), new(*MemberRepository)),
	wire.Bind(new(member.IUsageRepository), new(*UsageRepository)),
)
