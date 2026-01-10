package scheduler

import (
	"github.com/google/wire"
)

var WireSchedulerSet = wire.NewSet(
	NewSlotSchedulerService,
	NewCronScheduler,
)
