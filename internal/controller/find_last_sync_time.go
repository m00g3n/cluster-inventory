package controller

import (
	"time"
)

// TODO replace order
func findLastSyncTime(annotations map[string]string) (bool, time.Time) {
	_, found := annotations[lastKubeconfigSyncAnnotation]
	if !found {
		return false, time.Time{}
	}

	lastSyncTimeString := annotations[lastKubeconfigSyncAnnotation]
	lastSyncTime, err := time.Parse(time.RFC3339, lastSyncTimeString)
	if err != nil {
		return false, time.Time{}
	}
	return true, lastSyncTime
}
