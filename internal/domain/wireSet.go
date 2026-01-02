package domain

import (
	"github.com/google/wire"
	"member-pre/internal/domain/auth"
)

var WireDoMainSet = wire.NewSet(
	auth.NewAuthService,
)
