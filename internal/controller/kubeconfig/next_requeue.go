package kubeconfig

import (
	"time"
)

// nextRequeue - predicts rotationDuration for next requeue of GardenerCluster CR
func nextRequeue(now, lastSyncTime time.Time, rotationPeriod time.Duration, modifier float64) time.Duration {
	rotationPeriodWithModifier := modifier * rotationPeriod.Minutes()
	minutesToRequeue := rotationPeriodWithModifier - now.Sub(lastSyncTime).Minutes()
	if minutesToRequeue <= 0 {
		return time.Duration(rotationPeriodWithModifier * float64(time.Minute))
	}

	return time.Duration(minutesToRequeue * float64(time.Minute))
}
