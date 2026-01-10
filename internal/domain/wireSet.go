package domain

import (
	"github.com/google/wire"
	"member-pre/internal/domain/appointment"
	"member-pre/internal/domain/auth"
	"member-pre/internal/domain/member"
	"member-pre/internal/domain/slot"
	"member-pre/internal/domain/store"
	"member-pre/internal/domain/user"
)

var WireDoMainSet = wire.NewSet(
	user.NewUserService,
	auth.NewAuthService,
	store.NewStoreService,
	slot.NewTemplateService,
	slot.NewSlotService,
	member.NewMemberService,
	member.NewUsageService,
	appointment.NewAppointmentService,
	appointment.NewPaymentService,
	// 绑定接口
	wire.Bind(new(auth.IUserService), new(*user.UserService)),
	wire.Bind(new(appointment.IPaymentService), new(*appointment.PaymentService)),
)
