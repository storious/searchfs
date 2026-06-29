package namenode

import (
	"context"
	"time"

	"gdfs/internal/cluster"
)

func (n *NameNode) StartHealthChecker(ctx context.Context, interval time.Duration, cfg cluster.HealthConfig) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case now := <-ticker.C:
				n.registry.EvaluateHealth(now, cfg)
			}
		}
	}()
}
