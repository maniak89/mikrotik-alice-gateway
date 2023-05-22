package checker

import (
	"time"
)

type Config struct {
	SystemInfoCheckPeriod time.Duration `env:"SYSTEM_INFO_CHECK_PERIOD,default=5m"`
}
