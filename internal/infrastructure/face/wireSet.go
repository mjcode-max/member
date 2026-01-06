package face

import (
	"github.com/google/wire"
	"member-pre/internal/infrastructure/config"
	"member-pre/pkg/logger"
)

// NewHuaweiFRSServiceFromConfig 从配置创建华为云FRS服务
func NewHuaweiFRSServiceFromConfig(cfg *config.Config, log logger.Logger) FaceService {
	return NewHuaweiFRSService(&cfg.HuaweiFRS, log)
}

var WireFaceSet = wire.NewSet(
	NewHuaweiFRSServiceFromConfig,
)

