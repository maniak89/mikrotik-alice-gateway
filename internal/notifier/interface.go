package notifier

import (
	"context"

	"mikrotik-alice-gateway/internal/models/common"
)

type Notifier interface {
	NotifyHostChanged(ctx context.Context, router *common.Router, host *common.Host) error
}
